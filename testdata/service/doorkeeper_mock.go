package service

import (
	"context"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	"github.com/stretchr/testify/mock"
)

type DoorkeeperServiceMock struct {
	mock.Mock
}

func (service *DoorkeeperServiceMock) HashPassword(pass string) ([]byte, error) {
	args := service.Called(pass)
	return args.Get(0).([]byte), args.Error(1)
}

func (service *DoorkeeperServiceMock) VerifyAndParseToken(ctx context.Context, tk string) (string, error) {
	args := service.Called(ctx, tk)
	return args.String(0), args.Error(1)
}

func (service *DoorkeeperServiceMock) GenerateUserTokens(user domain.User) (vo.UserToken, error) {
	args := service.Called(user)
	return args.Get(0).(vo.UserToken), args.Error(1)
}

func (service *DoorkeeperServiceMock) GenerateSession(payload []byte) (string, error) {
	args := service.Called(payload)
	return args.String(0), args.Error(1)
}

func (service *DoorkeeperServiceMock) ParseSession(session string) ([]byte, error) {
	args := service.Called(session)
	return args.Get(0).([]byte), args.Error(1)
}

func (service *DoorkeeperServiceMock) VerifyAdminKey(adminKey string) error {
	args := service.Called(adminKey)
	return args.Error(0)
}

func (service *DoorkeeperServiceMock) VerifyAndParseRefreshToken(ctx context.Context, tk string) (string, error) {
	args := service.Called(ctx, tk)
	return args.String(0), args.Error(1)
}

func (service *DoorkeeperServiceMock) CompareHashAndPassword(ctx context.Context, password string, hash []byte) (bool, error) {
	args := service.Called(ctx, password, hash)
	return args.Bool(0), args.Error(1)
}
