package doorkeeper

import (
	"fmt"
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
		"userId": user.Id,
		"iat":    now,
		"eat":    now.Add(s.dk.Duration),
	}
	t := jwt.NewWithClaims(s.dk.GetSignMethod(), claims)
	res, err = t.SignedString(s.dk.GetPrivKey())
	if err != nil {
		return res, err
	}
	return res, nil
}
