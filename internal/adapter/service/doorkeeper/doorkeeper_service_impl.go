package doorkeeper

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
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

func (s *doorkeeperService) GenerateToken(user domain.User) (res string, err error) {
	now := time.Now().UTC()
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

func (s *doorkeeperService) VerifyAndParseToken(tk string) (string, error) {
	token, err := jwt.Parse(tk, func(t *jwt.Token) (interface{}, error) {
		switch s.dk.GetConcreteSignMethod() {
		case reflect.TypeOf(&jwt.SigningMethodRSA{}):
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("signing method invalid")
			}
		}
		return s.dk.GetPubKey(), nil
	})
	if err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("validate: invalid")
	}

	if err := s.verifyClaims(claims); err != nil {
		return "", err
	}

	return claims["userId"].(string), nil
}

func (s *doorkeeperService) verifyClaims(claims jwt.MapClaims) error {
	if time.Now().UTC().Unix() > int64(claims["eat"].(float64)) {
		return fmt.Errorf("token has expired")
	}

	if claims["iss"].(string) != s.dk.GetIssuer() {
		return fmt.Errorf("unrecognized issuer")
	}

	return nil
}
