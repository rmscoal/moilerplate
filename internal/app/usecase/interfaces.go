package usecase

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
)

type ICredentialUseCase interface {
	// Login handle user signin for first time user and generate pair of a jwts
	SignUp(ctx context.Context, user domain.User) (domain.User, error)
	// Login handle user login and generate pair of jwts
	Login(ctx context.Context, cred vo.UserCredential) (domain.User, error)
	// Authenticates authenticates user from the given jwt.
	Authenticate(ctx context.Context, token string) (domain.User, error)
	// Refresh validates refresh tokens and generates a new set of tokens.
	Refresh(ctx context.Context, token string) (domain.User, error)

	// AdminLogin handles for admin/developer login by only the admin secret. This
	// is used to access admin resource such as swagger documentation for now.
	AdminLogin(ctx context.Context, adminKey string) (vo.AdminSession, error)
	// AuthenticateAdmin authenticates the admin session.
	AuthenticateAdmin(ctx context.Context, session string) error
}

type IUserProfileUseCase interface {
	// RetrieveProfile fetches the user's full profile.
	RetrieveProfile(ctx context.Context, userID string) (domain.User, error)
	// ModifyEmailAddress modifies user emails
	ModifyEmailAddress(ctx context.Context, user domain.User) (domain.User, error)
}
