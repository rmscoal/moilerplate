package usecase

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type ICredentialUseCase interface {
	SignUp(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, cred vo.UserCredential) (domain.User, error)
}
