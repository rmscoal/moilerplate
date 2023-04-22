package service

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/stretchr/testify/mock"
)

type DoorkeeperServiceMock struct {
	mock.Mock
}

func (service *DoorkeeperServiceMock) HashPassword(pass string) string {
	args := service.Called(pass)
	return args.String(0)
}

func (service *DoorkeeperServiceMock) GenerateToken(user domain.User) (string, error) {
	args := service.Called(user)
	return args.String(0), args.Error(1)
}

func (service *DoorkeeperServiceMock) VerifyAndParseToken(tk string) (string, error) {
	args := service.Called(tk)
	return args.String(0), args.Error(1)
}
