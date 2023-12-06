package repo

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

type ICredentialRepo interface {
	IBaseRepo

	CreateNewUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByJti(ctx context.Context, jti string) (domain.User, error)
	SetNewUserToken(ctx context.Context, user domain.User) (vo.UserToken, error)
	UndoSetUserToken(ctx context.Context, jti string) error
	GetLatestUserTokenVersion(ctx context.Context, user domain.User) (int, error)
	DeleteUserTokenFamily(ctx context.Context, user domain.User) error
	RotateUserHashPassword(ctx context.Context, user domain.User) error
}
