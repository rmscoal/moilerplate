package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/stretchr/testify/mock"
)

type CredentialRepoMock struct {
	mock.Mock
}

func (repo *CredentialRepoMock) ValidateRepoState(ctx context.Context, user domain.User) error {
	args := repo.Called(ctx, user)
	return args.Error(0)
}

func (repo *CredentialRepoMock) CreateNewUser(ctx context.Context, user domain.User) (domain.User, error) {
	args := repo.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repo *CredentialRepoMock) GetUserByCredentials(ctx context.Context, cred vo.UserCredential) (domain.User, error) {
	args := repo.Called(ctx, cred)
	return args.Get(0).(domain.User), args.Error(1)
}
