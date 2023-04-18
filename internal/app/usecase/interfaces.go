package usecase

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

type ICredentialUseCase interface {
	SignUp(ctx context.Context, user domain.User) (domain.User, error)
}
