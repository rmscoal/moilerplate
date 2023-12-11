package usecase

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

type ICredentialUseCase interface {
	SignUp(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, cred vo.UserCredential) (domain.User, error)
	Authenticate(ctx context.Context, token string) (domain.User, error)
	Refresh(ctx context.Context, token string) (domain.User, error)

	// AdminLogin handles for admin/developer login by only the admin secret. This
	// is used to access admin resource such as swagger documentation for now.
	AdminLogin(ctx context.Context, adminKey string) (vo.AdminSession, error)
	// AuthenticateAdmin authenticates the admin session.
	AuthenticateAdmin(ctx context.Context, session string) error
}

type IUserProfileUseCase interface {
	RetrieveProfile(ctx context.Context, userID string) (domain.User, error)
	ModifyEmailAddress(ctx context.Context, user domain.User) (domain.User, error)
}
