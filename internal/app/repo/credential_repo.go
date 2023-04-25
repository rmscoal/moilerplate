package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type ICredentialRepo interface {
	ValidateRepoState(ctx context.Context, user domain.User) error
	CreateNewUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByCredentials(ctx context.Context, cred vo.UserCredential) (domain.User, error)
	GetUserByJti(ctx context.Context, jti string) (domain.User, error)
	SetNewUserToken(ctx context.Context, user domain.User) (vo.UserToken, error)
	UndoSetUserToken(ctx context.Context, jti string) error
	GetLatestUserTokenVersion(ctx context.Context, user domain.User) (int, error)
	DeleteUserTokenFamily(ctx context.Context, user domain.User) error
}
