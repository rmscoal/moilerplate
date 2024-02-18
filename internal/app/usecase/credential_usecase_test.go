package usecase

import (
	"context"
	"testing"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/testing/observability"
	mockrepo "github.com/rmscoal/moilerplate/testing/repo"
	mockservice "github.com/rmscoal/moilerplate/testing/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type CredentialUseCaseTestSuite struct {
	suite.Suite
	credRepo *mockrepo.CredentialRepoMock
	dkSvc    *mockservice.DoorkeeperServiceMock
}

func TestCredentialUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CredentialUseCaseTestSuite))
}

func (suite *CredentialUseCaseTestSuite) SetupSuite() {
	otel.SetTracerProvider(observability.NewTestingTracerProvider())
}

func (suite *CredentialUseCaseTestSuite) SetupTest() {
	suite.credRepo = new(mockrepo.CredentialRepoMock)
	suite.dkSvc = new(mockservice.DoorkeeperServiceMock)
}

func (suite *CredentialUseCaseTestSuite) TestSignUp_Success() {
	newUser := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	savedUser := newUser
	savedUser.ID = "some_id"

	suite.credRepo.On("CreateUser", context.Background(), newUser).Return(savedUser, nil)

	uc := NewCredentialUseCase(suite.credRepo, suite.dkSvc)
	result, err := uc.SignUp(context.Background(), newUser)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), result.ID)
}
