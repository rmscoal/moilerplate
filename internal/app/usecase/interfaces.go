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
	Login(ctx context.Context, cred domain.User) (vo.Token, error)
	// Authenticates authenticates user from the given jwt.
	Authenticate(ctx context.Context, token string) (domain.User, error)
	// Refresh validates refresh tokens and generates a new set of tokens.
	Refresh(ctx context.Context, token string) (vo.Token, error)
}
