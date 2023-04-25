package service

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/rmscoal/go-restful-monolith-boilerplate/pkg/doorkeeper"
)

type doorkeeperService struct {
	dk *doorkeeper.Doorkeeper
}

func NewDoorkeeperService(dk *doorkeeper.Doorkeeper) *doorkeeperService {
	return &doorkeeperService{dk}
}

func (s *doorkeeperService) HashPassword(pass string) string {
	h := s.dk.GetHasMethod().New()
	h.Write([]byte(pass))
	res := h.Sum([]byte(s.dk.GetSalt()))

	return fmt.Sprintf("%x", res)
}

func (s *doorkeeperService) GenerateUserTokens(user domain.User) (vo.UserToken, error) {
	var userToken vo.UserToken

	acct, err := s.GenerateAccessToken(user)
	if err != nil {
		return userToken, err
	}

	rt, err := s.GenerateRefreshToken(user)
	if err != nil {
		return userToken, err
	}

	userToken.AccesssToken = acct
	userToken.RefreshToken = rt

	return userToken, nil
}

func (s *doorkeeperService) GenerateAccessToken(user domain.User) (res string, err error) {
	now := user.Credential.Tokens.IssuedAt.UTC()
	claims := jwt.MapClaims{
		"iss":    s.dk.GetIssuer(),
		"eat":    now.Add(s.dk.Duration).Unix(),
		"iat":    now.Unix(),
		"userId": user.Id,
		"nbf":    now.Unix(),
	}

	res, err = jwt.NewWithClaims(s.dk.GetSignMethod(), claims).SignedString(s.dk.GetPrivKey())
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *doorkeeperService) GenerateRefreshToken(user domain.User) (rt string, err error) {
	now := user.Credential.Tokens.IssuedAt.UTC()

	claims := jwt.MapClaims{
		"iss": s.dk.GetIssuer(),
		"eat": now.Add(time.Hour /* use: s.dk.RefreshDuration */).Unix(),
		"iat": now.Unix(),
		"jti": user.Credential.Tokens.TokenID,
		"nbf": now.Unix(),
	}

	rt, err = jwt.NewWithClaims(s.dk.GetSignMethod(), claims).SignedString(s.dk.GetPrivKey())
	if err != nil {
		return rt, err
	}

	return rt, nil
}

func (s *doorkeeperService) VerifyAndParseToken(ctx context.Context, tk string) (string, error) {
	claims, err := s.verifyAndGetClaims(tk)
	if err != nil {
		return "", err
	}

	if err := s.verifyClaims(ctx, claims, "userId"); err != nil {
		return "", err
	}

	return claims["userId"].(string), nil
}

func (s *doorkeeperService) VerifyAndParseRefreshToken(ctx context.Context, tk string) (string, error) {
	claims, err := s.verifyAndGetClaims(tk)
	if err != nil {
		return "", err
	}

	if err := s.verifyClaims(ctx, claims, "jti"); err != nil {
		return "", err
	}

	return claims["jti"].(string), nil
}

func (s *doorkeeperService) verifyAndGetClaims(tk string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tk, func(t *jwt.Token) (interface{}, error) {
		switch s.dk.GetConcreteSignMethod() {
		case doorkeeper.RSA_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("signing method invalid")
			}
		case doorkeeper.RSAPSS_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodRSAPSS); !ok {
				return nil, fmt.Errorf("signing method invalid")
			}
		case doorkeeper.HMAC_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method invalid")
			}
		case doorkeeper.ECDSA_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("signing method invalid")
			}
		case doorkeeper.EdDSA_SIGN_METHOD_TYPE:
			if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("signing method invalid")
			}
		}
		return s.dk.GetPubKey(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}

func (s *doorkeeperService) verifyClaims(ctx context.Context, claims jwt.MapClaims, expectedKeys ...any) error {
	keys := []any{"iss", "iat", "eat", "nbf"}
	keys = append(keys, []any(expectedKeys)...)

	if err := s.validateKeys(ctx, claims, keys...); err != nil {
		return err
	}

	now := time.Now().UTC()

	if _, ok := claims["iss"].(string); !ok {
		return fmt.Errorf("invalid token claims")
	}
	if _, ok := claims["eat"].(float64); !ok {
		return fmt.Errorf("invalid token claims")
	}
	if _, ok := claims["nbf"].(float64); !ok {
		return fmt.Errorf("invalid token claims")
	}

	if now.Unix() > int64(claims["eat"].(float64)) {
		return fmt.Errorf("token has expired")
	}

	if int64(claims["nbf"].(float64)) > now.Unix() {
		return fmt.Errorf("invalid token claims: nbf > now")
	}

	if claims["iss"].(string) != s.dk.GetIssuer() {
		return fmt.Errorf("unrecognized issuer")
	}

	return nil
}

func (s *doorkeeperService) validateKeys(ctx context.Context, obj map[string]any, args ...any) error {
	keys := make([]string, len(obj))

	index := 0
	for k := range obj {
		keys[index] = k
		index++
	}

	return validation.ValidateWithContext(ctx, keys, validation.Each(validation.Required, validation.In(args...)))
}
