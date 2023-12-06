package service

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

type IDoorkeeperService interface {
	// GenerateSession generates a sessions by the given payload. The payload
	// will be aes encrypted and base64 encoded.
	GenerateSession(payload []byte) (string, error)

	// ParseSession decode the session and then decrypts it giving back the
	// original payload.
	ParseSession(session string) ([]byte, error)

	// VerifyAdminKey checks whether the key given matches with the original admin key.
	VerifyAdminKey(adminKey string) error

	HashPassword(pass string) ([]byte, error)
	CompareHashAndPassword(ctx context.Context, password string, hash []byte) (bool, error)
	GenerateUserTokens(user domain.User) (vo.UserToken, error)
	VerifyAndParseToken(ctx context.Context, tk string) (string, error)
	VerifyAndParseRefreshToken(ctx context.Context, tk string) (string, error)
}
