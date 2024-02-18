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
	ctx := context.Background()
	newUser := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	encodedUser := newUser
	encodedUser.Password = "hashed string"
	suite.dkSvc.On("HashAndEncodeStringWithSalt", ctx, newUser.Password, newUser.Username).Return("hashed string")

	savedUser := encodedUser
	savedUser.ID = "some id"
	suite.credRepo.On("CreateUser", ctx, encodedUser).Return(savedUser, nil)

	uc := NewCredentialUseCase(suite.credRepo, suite.dkSvc)
	result, err := uc.SignUp(ctx, newUser)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), result.ID)
	assert.Equal(suite.T(), savedUser.Password, result.Password)
}

func (suite *CredentialUseCaseTestSuite) TestSignup_Fail_Validation() {
	ctx := context.Background()
	newUser := domain.User{
		Password: "verystrongpassword",
	}

	uc := NewCredentialUseCase(suite.credRepo, suite.dkSvc)
	result, err := uc.SignUp(ctx, newUser)
	assert.ErrorContains(suite.T(), err, ErrUnprocessableEntity.Error())
	assert.Equal(suite.T(), newUser, result)
}

func (suite *CredentialUseCaseTestSuite) TestSignup_Fail_CreateUserRepo() {
	ctx := context.Background()
	newUser := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	encodedUser := newUser
	encodedUser.Password = "hashed string"
	suite.dkSvc.On("HashAndEncodeStringWithSalt", ctx, newUser.Password, newUser.Username).Return("hashed string")

	suite.credRepo.On("CreateUser", ctx, encodedUser).Return(encodedUser, ErrUnexpected)

	uc := NewCredentialUseCase(suite.credRepo, suite.dkSvc)
	result, err := uc.SignUp(ctx, newUser)
	assert.ErrorContains(suite.T(), err, ErrUnexpected.Error())
	assert.Equal(suite.T(), encodedUser, result)
}
