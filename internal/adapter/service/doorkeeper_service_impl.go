package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"github.com/rmscoal/moilerplate/internal/utils"
	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"golang.org/x/crypto/pbkdf2"
)

var (
	MinSaltLength int64 = 1 << 5
	MaxSaltLength int64 = 1 << 6
	NumWorkers    int   = 10
)

var ErrPasswordMismatch = errors.New("timeout exceeded due to password mismatch")

type doorkeeperService struct {
	dk *doorkeeper.Doorkeeper
}

func NewDoorkeeperService(dk *doorkeeper.Doorkeeper) *doorkeeperService {
	return &doorkeeperService{dk}
}

/*
---------- Admin Section ----------
*/

func (s *doorkeeperService) VerifyAdminKey(adminKey string) error {
	// Since the adminKey is hashed using sha256, we have to compare
	// by also using sha256.
	sha := sha256.Sum256(utils.ConvertStringToByteSlice(adminKey))
	if subtle.ConstantTimeCompare(
		sha[:],
		utils.ConvertStringToByteSlice(s.dk.GetAdminKey()),
	) == 0 {
		return fmt.Errorf("invalid admin key")
	}

	return nil
}

// GenerateSession generates a sessions by the given payload. The payload
// will be aes encrypted and base64 encoded.
func (s *doorkeeperService) GenerateSession(payload []byte) (string, error) {
	encrypted, err := aesGCMEncrypt(utils.ConvertStringToByteSlice(s.dk.GetAdminKey()), payload)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// ParseSession decode the session and then decrypts it giving back the
// original payload.
func (s *doorkeeperService) ParseSession(session string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(session)
	if err != nil {
		return nil, errors.New("failed to decode session")
	}

	plaintext, err := aesGCMDecrypt(utils.ConvertStringToByteSlice(s.dk.GetAdminKey()), ciphertext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

/*
---------- Hashing Section ----------
*/

func (s *doorkeeperService) HashPassword(pass string) ([]byte, error) {
	saltLength, err := rand.Int(rand.Reader, big.NewInt(MaxSaltLength-MinSaltLength))
	if err != nil {
		return nil, err
	}

	skipper, _ := rand.Int(rand.Reader, big.NewInt(2))

	salt := make([]byte, saltLength.Int64()+MinSaltLength)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	mixture := pbkdf2.Key(utils.ConvertStringToByteSlice(pass), salt, s.dk.GetHashIter(), s.dk.GetHashKeyLen(), s.dk.GetHasherFunc())

	for i, j := 0, 0; j < len(salt); i, j = i+int(skipper.Int64())+1, j+1 {
		mixture = utils.InsertAt(mixture, i, salt[j])
	}
	return mixture, nil
}

func (s *doorkeeperService) CompareHashAndPassword(ctx context.Context, password string, hash []byte) (bool, error) {
	reportChannel := make(chan bool)

	for i := MinSaltLength; i <= MaxSaltLength; i++ {
		go s.compareWorker(ctx, password, hash, i, reportChannel)
	}

	select {
	case <-ctx.Done():
		return false, ErrPasswordMismatch
	case <-reportChannel:
		return true, nil
	}
}

func (s *doorkeeperService) compareWorker(
	ctx context.Context, password string, hashToExtract []byte,
	lengthOfSalt int64, reportChannel chan<- bool,
) {
	for i := 0; i < 2; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			saltPrediction, hashPrediction := s.extractFromMixtures(hashToExtract, i, lengthOfSalt)
			if s.compareHashes(password, saltPrediction, hashPrediction) {
				reportChannel <- true
				return
			}
		}
	}
}

func (s *doorkeeperService) extractFromMixtures(hashToExtract []byte, skipper int, lengthOfSalt int64) ([]byte, []byte) {
	skipperIdx := 0
	saltCollected := make([]byte, 0, lengthOfSalt)
	hashCollected := make([]byte, 0, s.dk.GetHashKeyLen())
	for i := 0; i < len(hashToExtract); i++ {
		if i == skipperIdx && len(saltCollected) < int(lengthOfSalt) {
			saltCollected = append(saltCollected, hashToExtract[i])
			skipperIdx += skipper + 1
		} else {
			hashCollected = append(hashCollected, hashToExtract[i])
		}
	}
	return saltCollected, hashCollected
}

func (s *doorkeeperService) compareHashes(password string, salt, hashToCompare []byte) bool {
	hash := pbkdf2.Key([]byte(password), salt, s.dk.GetHashIter(), s.dk.GetHashKeyLen(), s.dk.GetHasherFunc())

	return subtle.ConstantTimeCompare(hash, hashToCompare) == 1
}

/*
---------- JWT Section ----------
*/

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
		"iss": s.dk.GetIssuer(),
		"eat": now.Add(s.dk.AccessDuration).Unix(),
		"iat": now.Unix(),
		"sub": user.Id,
		"nbf": now.Unix(),
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
		"eat": now.Add(s.dk.RefreshDuration).Unix(),
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

	if err := s.verifyClaims(ctx, claims, "sub"); err != nil {
		return "", err
	}

	return claims["sub"].(string), nil
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

	return validation.ValidateWithContext(ctx, keys, validation.Each(validation.Required, validation.In(args...).
		Error("does not contained required claim")))
}

/*
---------- Encryption Section ----------
*/

func aesGCMEncrypt(key []byte, payload []byte) ([]byte, error) {
	gcm, err := initAES_GCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("failed to read random")
	}

	return gcm.Seal(nonce, nonce, payload, nil), nil
}

func aesGCMDecrypt(key []byte, ciphertext []byte) ([]byte, error) {
	gcm, err := initAES_GCM(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("invalid size")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func initAES_GCM(key []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to initialize encryption")
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, errors.New("failed to authenticate admin session")
	}

	return gcm, nil
}
