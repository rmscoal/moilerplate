package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DoorkeeperServiceMock struct {
	mock.Mock
}

func (service *DoorkeeperServiceMock) HashAndEncodeStringWithSalt(ctx context.Context, str, slt string) string {
	args := service.Called(ctx, str, slt)
	return args.String(0)
}
