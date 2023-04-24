package service

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

type IDoorkeeperService interface {
	HashPassword(pass string) string
	GenerateToken(user domain.User) (string, error)
	VerifyAndParseToken(ctx context.Context, tk string) (string, error)
}
