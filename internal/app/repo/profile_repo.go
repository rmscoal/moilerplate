package repo

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
)

type IUserProfileRepo interface {
	IBaseRepo

	GetUserProfile(ctx context.Context, id string) (domain.User, error)
	SaveUserEmails(ctx context.Context, user domain.User) (domain.User, error)
}
