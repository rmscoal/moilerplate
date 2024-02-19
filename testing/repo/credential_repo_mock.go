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

func (repo *CredentialRepoMock) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	args := repo.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repo *CredentialRepoMock) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	args := repo.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}
