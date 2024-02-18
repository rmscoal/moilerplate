package repo

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
)

type ICredentialRepo interface {
	IBaseRepo

	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}
