package repo

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/stretchr/testify/mock"
)

type CredentialRepoMock struct {
	BaseRepoMock
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

func (repo *CredentialRepoMock) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	args := repo.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repo *CredentialRepoMock) GetUserByJti(ctx context.Context, jti string) (domain.User, error) {
	args := repo.Called(ctx, jti)
	return args.Get(0).(domain.User), args.Error(1)
}

func (repo *CredentialRepoMock) SetNewUserToken(ctx context.Context, user domain.User) (vo.UserToken, error) {
	args := repo.Called(ctx, user)
	return args.Get(0).(vo.UserToken), args.Error(1)
}

func (repo *CredentialRepoMock) UndoSetUserToken(ctx context.Context, jti string) error {
	args := repo.Called(ctx, jti)
	return args.Error(0)
}

func (repo *CredentialRepoMock) GetLatestUserTokenVersion(ctx context.Context, user domain.User) (int, error) {
	args := repo.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (repo *CredentialRepoMock) DeleteUserTokenFamily(ctx context.Context, user domain.User) error {
	args := repo.Called(ctx, user)
	return args.Error(0)
}

func (repo *CredentialRepoMock) RotateUserHashPassword(ctx context.Context, user domain.User) error {
	args := repo.Called(ctx, user)
	return args.Error(0)
}
