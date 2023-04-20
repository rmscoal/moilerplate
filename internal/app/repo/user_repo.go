package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type IUserRepo interface {
	ValidateRepoState(ctx context.Context, user domain.User) error
	CreateNewUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByCredentials(ctx context.Context, cred vo.UserCredential) (domain.User, error)
}
