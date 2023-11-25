package repo

import (
	"github.com/stretchr/testify/mock"
)

type BaseRepoMock struct {
	mock.Mock
}

func (repo *BaseRepoMock) TranslateError(err error) error {
	args := repo.Called(err)
	return args.Error(0)
}
