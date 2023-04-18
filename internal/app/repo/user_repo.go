package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

type IUserRepo interface {
	CreateNewUser(ctx context.Context, user domain.User) (domain.User, error)
}
