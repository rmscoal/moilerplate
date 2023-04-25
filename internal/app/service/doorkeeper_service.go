package service

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type IDoorkeeperService interface {
	HashPassword(pass string) string
	GenerateUserTokens(user domain.User) (vo.UserToken, error)
	VerifyAndParseToken(ctx context.Context, tk string) (string, error)
	VerifyAndParseRefreshToken(ctx context.Context, tk string) (string, error)
}
