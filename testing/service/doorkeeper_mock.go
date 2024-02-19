package service

import (
	"context"

	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"github.com/stretchr/testify/mock"
)

type DoorkeeperServiceMock struct {
	mock.Mock
}

func (service *DoorkeeperServiceMock) HashAndEncodeStringWithSalt(ctx context.Context, str, slt string) string {
	args := service.Called(ctx, str, slt)
	return args.String(0)
}

func (service *DoorkeeperServiceMock) ComparePasswords(ctx context.Context, hashAndEncodedPass, passToCheck, salt string) (bool, error) {
	args := service.Called(ctx, hashAndEncodedPass, passToCheck, salt)
	return args.Get(0).(bool), args.Error(1)
}

func (service *DoorkeeperServiceMock) GenerateTokens(ctx context.Context, subject string, prevJTI *string) (vo.Token, error) {
	args := service.Called(ctx, subject, prevJTI)
	return args.Get(0).(vo.Token), args.Error(1)
}
