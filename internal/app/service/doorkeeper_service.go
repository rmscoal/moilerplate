package service

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

type IDoorkeeperService interface {
	HashAndEncodeStringWithSalt(ctx context.Context, str, slt string) string
	ComparePasswords(ctx context.Context, hashAndEncodedPass, passToCheck, salt string) (bool, error)
	GenerateTokens(ctx context.Context, subject string, prevJTI *string) (vo.Token, error)
	ValidateAccessToken(ctx context.Context, token string) (string, error)
	ValidateRefreshToken(ctx context.Context, token string) (vo.Token, error)
}
