package usecase

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
)

type ICredentialUseCase interface {
	SignUp(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, cred vo.UserCredential) (domain.User, error)
	Authorize(ctx context.Context, token string) (domain.User, error)
	AuthAdmin(ctx context.Context, adminKey string) error
	Refresh(ctx context.Context, token string) (domain.User, error)
}

type IUserProfileUseCase interface {
	RetrieveProfile(ctx context.Context, userID string) (domain.User, error)
	ModifyEmailAddress(ctx context.Context, user domain.User) (domain.User, error)
}
