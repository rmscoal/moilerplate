package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
)

type IUserProfileRepo interface {
	GetUserProfile(ctx context.Context, id string) (domain.User, error)
	SaveUserEmails(ctx context.Context, user domain.User) (domain.User, error)
}
