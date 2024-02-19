package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"github.com/rmscoal/moilerplate/testing/observability"
	mockrepo "github.com/rmscoal/moilerplate/testing/repo"
	mockservice "github.com/rmscoal/moilerplate/testing/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type credentialUseCaseTestSuite struct {
	suite.Suite
	credRepo *mockrepo.CredentialRepoMock
	dkSvc    *mockservice.DoorkeeperServiceMock
	uc       ICredentialUseCase
}

func TestCredentialUseCaseSuite(t *testing.T) {
	suite.Run(t, new(credentialUseCaseTestSuite))
}

func (suite *credentialUseCaseTestSuite) SetupSuite() {
	otel.SetTracerProvider(observability.NewTestingTracerProvider())
}

func (suite *credentialUseCaseTestSuite) SetupTest() {
	suite.credRepo = new(mockrepo.CredentialRepoMock)
	suite.dkSvc = new(mockservice.DoorkeeperServiceMock)

	suite.uc = NewCredentialUseCase(suite.credRepo, suite.dkSvc)
}

func (suite *credentialUseCaseTestSuite) TestSignUp_Success() {
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

	result, err := suite.uc.SignUp(ctx, newUser)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), result.ID)
	assert.Equal(suite.T(), savedUser.Password, result.Password)
}

func (suite *credentialUseCaseTestSuite) TestSignup_Fail_Validation() {
	ctx := context.Background()
	newUser := domain.User{
		Password: "verystrongpassword",
	}

	result, err := suite.uc.SignUp(ctx, newUser)
	assert.ErrorContains(suite.T(), err, ErrUnprocessableEntity.Error())
	assert.Equal(suite.T(), newUser, result)
}

func (suite *credentialUseCaseTestSuite) TestSignup_Fail_CreateUserRepo() {
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

	result, err := suite.uc.SignUp(ctx, newUser)
	assert.ErrorContains(suite.T(), err, ErrUnexpected.Error())
	assert.Equal(suite.T(), encodedUser, result)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Success() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "verystrongpassword",
	}
	user := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	suite.credRepo.On("GetUserByUsername", ctx, cred.Username).Return(user, nil)
	suite.dkSvc.On("ComparePasswords", ctx, user.Password, cred.Password, user.Username).Return(true, nil)
	suite.dkSvc.On("GenerateTokens", ctx, user.ID, mock.Anything).Return(vo.Token{AccessToken: "AT", RefreshToken: "RT"}, nil)

	token, err := suite.uc.Login(ctx, cred)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Fail_Validation() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "",
	}

	token, err := suite.uc.Login(ctx, cred)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Fail_GetUserByUsername_NotFound() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "verystrongpassword",
	}

	suite.credRepo.On("GetUserByUsername", ctx, cred.Username).Return(domain.User{}, ErrNotFound)

	token, err := suite.uc.Login(ctx, cred)
	assert.ErrorContains(suite.T(), err, ErrUnauthorized.Error())
	assert.Empty(suite.T(), token)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Fail_GetUserByUsername_Unknown() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "verystrongpassword",
	}

	suite.credRepo.On("GetUserByUsername", ctx, cred.Username).Return(domain.User{}, ErrUnexpected)

	token, err := suite.uc.Login(ctx, cred)
	assert.ErrorContains(suite.T(), err, ErrUnauthorized.Error())
	assert.Empty(suite.T(), token)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Fail_ComparePasswords_Not_Match() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "verystrongpassword",
	}
	user := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	suite.credRepo.On("GetUserByUsername", ctx, cred.Username).Return(user, nil)
	suite.dkSvc.On("ComparePasswords", ctx, user.Password, cred.Password, user.Username).Return(false, nil)

	token, err := suite.uc.Login(ctx, cred)
	assert.ErrorContains(suite.T(), err, ErrUnauthorized.Error())
	assert.Empty(suite.T(), token)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Fail_ComparePasswords_Decode() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "verystrongpassword",
	}
	user := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	suite.credRepo.On("GetUserByUsername", ctx, cred.Username).Return(user, nil)
	suite.dkSvc.On("ComparePasswords", ctx, user.Password, cred.Password, user.Username).Return(false, errors.New("decode error"))

	token, err := suite.uc.Login(ctx, cred)
	assert.ErrorContains(suite.T(), err, ErrUnexpected.Error())
	assert.Empty(suite.T(), token)
}

func (suite *credentialUseCaseTestSuite) TestLogin_Fail_GenerateTokens() {
	ctx := context.Background()
	cred := domain.User{
		Username: "rmscoal",
		Password: "verystrongpassword",
	}
	user := domain.User{
		Name:        "Rifky Satyana",
		Username:    "rmscoal",
		Email:       "rmscoaldev@gmail.com",
		PhoneNumber: "6281234274916",
		Password:    "verystrongpassword",
	}

	suite.credRepo.On("GetUserByUsername", ctx, cred.Username).Return(user, nil)
	suite.dkSvc.On("ComparePasswords", ctx, user.Password, cred.Password, user.Username).Return(true, nil)
	suite.dkSvc.On("GenerateTokens", ctx, user.ID, mock.Anything).Return(vo.Token{}, errors.New("unknown"))

	token, err := suite.uc.Login(ctx, cred)
	assert.ErrorContains(suite.T(), err, ErrUnexpected.Error())
	assert.Empty(suite.T(), token)
}
