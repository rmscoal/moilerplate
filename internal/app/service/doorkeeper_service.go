package service

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type IDoorkeeperService interface {
	HashPassword(pass string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, password string, hash []byte) (bool, error)
	GenerateUserTokens(user domain.User) (vo.UserToken, error)
	VerifyAdminKey(adminKey string) error
	VerifyAndParseToken(ctx context.Context, tk string) (string, error)
	VerifyAndParseRefreshToken(ctx context.Context, tk string) (string, error)
}
