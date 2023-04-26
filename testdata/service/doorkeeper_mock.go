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

func (service *DoorkeeperServiceMock) HashPassword(pass string) string {
	args := service.Called(pass)
	return args.String(0)
}

func (service *DoorkeeperServiceMock) VerifyAndParseToken(ctx context.Context, tk string) (string, error) {
	args := service.Called(ctx, tk)
	return args.String(0), args.Error(1)
}

func (service *DoorkeeperServiceMock) GenerateUserTokens(user domain.User) (vo.UserToken, error) {
	args := service.Called(user)
	return args.Get(0).(vo.UserToken), args.Error(1)
}

func (service *DoorkeeperServiceMock) VerifyAndParseRefreshToken(ctx context.Context, tk string) (string, error) {
	args := service.Called(ctx, tk)
	return args.String(0), args.Error(1)
}
