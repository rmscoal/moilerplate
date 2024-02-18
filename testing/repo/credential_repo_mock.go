package repo

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CredentialRepoMock struct {
	BaseRepoMock
	mock.Mock
}

func (repo *CredentialRepoMock) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	args := repo.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}
