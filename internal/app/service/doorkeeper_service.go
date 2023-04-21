package service

import "github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"

type IDoorkeeperService interface {
	HashPassword(pass string) string
	GenerateToken(user domain.User) (string, error)
	VerifyAndParseToken(tk string) (string, error)
}
