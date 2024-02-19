package service

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/pbkdf2"
)

type doorkeeperService struct {
	db     *sql.DB // For storing/retrieving access versioning table
	dk     *doorkeeper.Doorkeeper
	tracer trace.Tracer
}

func NewDoorkeeperService(dk *doorkeeper.Doorkeeper, db *sql.DB) *doorkeeperService {
	return &doorkeeperService{dk: dk, db: db, tracer: otel.Tracer("doorkeeper-service")}
}

func (service *doorkeeperService) HashAndEncodeStringWithSalt(ctx context.Context, str, slt string) string {
	_, span := service.tracer.Start(ctx, "service.HashAndEncodeStringWithSalt")
	span.End()

	salt := sha1.Sum([]byte(slt))

	return base64.StdEncoding.EncodeToString(
		pbkdf2.Key(
			[]byte(str),
			salt[:],
			service.dk.GetHashIter(),
			service.dk.GetHashKeyLen(),
			service.dk.GetHasherFunc(),
		),
	)
}

/*
---------- JWT Section ----------
*/

var (
	ErrGenerateJTI               error = errors.New("failed generating jti")
	ErrCreateToken               error = errors.New("failed generating new token")
	ErrStoreJWTInRedis           error = errors.New("unable to store token in Redis")
	ErrInvalidSignMethod         error = errors.New("signing method invalid")
	ErrInvalidClaims             error = errors.New("invalid claims")
	ErrTokenExpiredOrInvalidated error = errors.New("token expired or invalidated")

	JWTRedisSubjectField        string = "subject"
	JWTRedisAccessTokenJTIField string = "atJTI"
)

func (service *doorkeeperService) GenerateTokens(ctx context.Context, subject string, prevJTI *string) (vo.Token, error) {
	ctx, span := service.tracer.Start(ctx, "service.GenerateTokens")
	defer span.End()

	var token vo.Token

	accessTokenClaims := jwt.MapClaims{
		"iss": service.dk.GetIssuer(),
		"iat": time.Now().Unix(),
		"sub": subject,
		"exp": time.Now().Add(service.dk.GetJWTAccessDuration()).Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		"iss": service.dk.GetIssuer(),
		"iat": time.Now().Unix(),
		"jti": uuid.NewString(),
		"exp": time.Now().Add(service.dk.GetJWTRefreshDuration()).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(service.dk.GetSignMethod(), accessTokenClaims).SignedString(service.dk.GetPrivKey())
	if err != nil {
		span.SetStatus(codes.Error, "failed to create access token")
		span.RecordError(err)
		return token, err
	}

	refreshToken, err := jwt.NewWithClaims(service.dk.GetSignMethod(), refreshTokenClaims).SignedString(service.dk.GetPrivKey())
	if err != nil {
		span.SetStatus(codes.Error, "failed to create refresh token")
		span.RecordError(err)
		return token, err
	}

	query := `
		INSERT INTO access_versionings (jti, parent_id, user_id, version) 
		VALUES (
			$1, $2, $3, (
			SELECT
				CASE $2 = '' OR $2 IS NULL
					WHEN TRUE THEN 1
					ELSE
						CASE (SELECT COUNT(*) FROM access_versionings WHERE jti = $2)
							WHEN 0 THEN 1
							ELSE (SELECT access_versionings.version + 1 FROM access_versionings WHERE jti = $2)
						END
				END AS version
			)
		)
	`

	span.SetAttributes(semconv.DBSystemPostgreSQL, semconv.DBStatementKey.String(query))
	if _, err = service.db.ExecContext(ctx, query, refreshTokenClaims["jti"], prevJTI, subject); err != nil {
		span.SetStatus(codes.Error, "failed to create refresh token")
		span.RecordError(err)
		return token, err
	}

	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	return token, nil
}

func (service *doorkeeperService) ValidateAccessToken(ctx context.Context, at string) (userID string, err error) {
	_, span := service.tracer.Start(ctx, "service.ValidateAccessToken")
	defer span.End()

	jwt, err := service.parseToken(at)
	if err != nil {
		return userID, err
	}

	return jwt.Claims.GetSubject()
}

func (service *doorkeeperService) ValidateRefreshToken(ctx context.Context, rt string) (token vo.Token, err error) {
	ctx, span := service.tracer.Start(ctx, "service.ValidateRefreshToken")
	defer span.End()

	tk, err := service.parseToken(rt)
	if err != nil {
		return token, err
	}

	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return token, ErrInvalidClaims
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return token, ErrInvalidClaims
	}

	query := `
		SELECT
			 av1.user_id
		FROM access_versionings av1
		WHERE
			jti = $1 AND 
			(SELECT av2.version FROM access_versionings av2 WHERE av2.user_id = av1.user_id ORDER BY version DESC LIMIT 1) = av1.version;
	`
	span.SetAttributes(semconv.DBSystemPostgreSQL, semconv.DBStatementKey.String(query))

	var userID string
	err = service.db.QueryRowContext(ctx, query, jti).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			service.InvalidateToken(ctx, jti)
			return token, ErrTokenExpiredOrInvalidated
		}

		return token, err
	}

	return service.GenerateTokens(ctx, userID, &jti)
}

func (service *doorkeeperService) InvalidateToken(ctx context.Context, jti string) error {
	ctx, span := service.tracer.Start(ctx, "service.InvalidateToken")
	defer span.End()

	query := `DELETE FROM access_versionings WHERE user_id = (SELECT user_id FROM access_versionings WHERE jti = $1)`
	span.SetAttributes(semconv.DBSystemPostgreSQL, semconv.DBStatementKey.String(query))

	_, err := service.db.ExecContext(ctx, query, jti)
	return err
}

func (service *doorkeeperService) parseToken(tk string) (*jwt.Token, error) {
	return jwt.Parse(tk, func(t *jwt.Token) (interface{}, error) {
		switch service.dk.GetConcreteSignMethod() {
		case doorkeeper.RSA_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, ErrInvalidSignMethod
			}
		case doorkeeper.RSAPSS_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodRSAPSS); !ok {
				return nil, ErrInvalidSignMethod
			}
		case doorkeeper.HMAC_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSignMethod
			}
		case doorkeeper.ECDSA_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, ErrInvalidSignMethod
			}
		case doorkeeper.EdDSA_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, ErrInvalidSignMethod
			}
		}
		return service.dk.GetPubKey(), nil
	},
		jwt.WithIssuer(service.dk.GetIssuer()),
		jwt.WithStrictDecoding(),
		jwt.WithLeeway(5*time.Minute),
	)
}
