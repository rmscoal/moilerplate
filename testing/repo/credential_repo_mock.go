package repo

import (
	"github.com/stretchr/testify/mock"
)

type CredentialRepoMock struct {
	BaseRepoMock
	mock.Mock
}
