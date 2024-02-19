package repo

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
)

type ICredentialRepo interface {
	IBaseRepo

	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}
