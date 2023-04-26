package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
	mockrepo "github.com/rmscoal/go-restful-monolith-boilerplate/testdata/repo"
	mockservice "github.com/rmscoal/go-restful-monolith-boilerplate/testdata/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

/*
Scenario:

A new user registers via its,
1. firstname,
2. lastname,
3. username,
4. password,
5. email, and
6. phone number

The creation of a new user is valid if:
a. the firstname is not empty and is around 3 - 20 bytes long
b. the lastname is not empty and is around 3 - 25 bytes long
c. the username is not empty and is around 3 - 20 bytes long
d. the email is not empty and is a valid email checked by regex
e. the password is not empty
f. the phone number is not empty and is a valid indonesian phone number

A user gains authorization by
1. input username and
2. input password

Once successful, an access_token and refresh_token is created
*/

var (
	// Context for testing
	TEST_CTX = context.Background()

	VALID_USER_DOMAIN = domain.User{
		FirstName: "RMSCOAL",
		LastName:  "RMSCOAL",
		Emails: []vo.UserEmail{
			{
				Email:     "rmscoal@test.com",
				IsPrimary: true,
			},
		},
		PhoneNumber: "0812345699",
		Credential: vo.UserCredential{
			Username: "RMSCOAL",
			Password: "PASSWORD",
		},
	}
	VALID_USER_DOMAIN_WITH_ID = domain.User{
		Id:        "UNIQUE_ID",
		FirstName: "RMSCOAL",
		LastName:  "RMSCOAL",
		Emails: []vo.UserEmail{
			{
				Email:     "rmscoal@test.com",
				IsPrimary: true,
			},
		},
		PhoneNumber: "0812345699",
		Credential: vo.UserCredential{
			Username: "RMSCOAL",
			Password: "PASSWORD",
		},
	}
	// Use this variable for testing unsuccessful cases.
	// It makes readibility better and understanable.
	INVALID_USER_DOMAIN = VALID_USER_DOMAIN

	VALID_USER_TOKENS_ONLY_ID = vo.UserToken{
		TokenID:  "TOKEN_ID",
		Issued:   false,
		IssuedAt: time.Now(),
	}

	// User Tokens
	VALID_USER_TOKENS = vo.UserToken{
		TokenID:      "TOKEN_ID",
		AccesssToken: "ACCESS_TOKEN",
		RefreshToken: "REFRESH_TOKEN",
		Version:      1,
		Issued:       false,
		IssuedAt:     time.Now(),
	}
)

type CredentialUseCaseTestSuite struct {
	suite.Suite
	repo    *mockrepo.CredentialRepoMock
	service *mockservice.DoorkeeperServiceMock
}

func (suite *CredentialUseCaseTestSuite) SetupTest() {
	suite.repo = new(mockrepo.CredentialRepoMock)
	suite.service = new(mockservice.DoorkeeperServiceMock)
}

func TestCredentialUseCase(t *testing.T) {
	suite.Run(t, new(CredentialUseCaseTestSuite))
}

func (suite *CredentialUseCaseTestSuite) TestSignup() {
	suite.Run("Successful Signup", func() {
		test := VALID_USER_DOMAIN_WITH_ID
		test.Credential.Password = "HASHED_PASSWORD"

		// Assumes that it passes all validity repo state checks
		suite.repo.On("ValidateRepoState", TEST_CTX, VALID_USER_DOMAIN).Return(nil).Once()
		// Assumes that it passes constraint checks while persisting record
		suite.repo.On("CreateNewUser", TEST_CTX, mock.AnythingOfType("domain.User")).Return(test, nil).Once()
		// Assumes preparing generating refresh token family done successfully
		suite.repo.On("SetNewUserToken", TEST_CTX, mock.AnythingOfType("domain.User")).Return(VALID_USER_TOKENS_ONLY_ID, nil).Once()
		// Assumes HashPassword call returns the hashed password
		suite.service.On("HashPassword", mock.AnythingOfType("string")).Return("HASHED_PASSWORD").Once()
		// Assumes HashPassword call successfully return the hashed password
		suite.service.On("GenerateUserTokens", mock.AnythingOfType("domain.User")).Return(VALID_USER_TOKENS, nil).Once()

		test.Credential.Tokens = VALID_USER_TOKENS

		// Start sign up test
		uc := NewCredentialUseCase(suite.repo, suite.service)
		user, err := uc.SignUp(TEST_CTX, VALID_USER_DOMAIN)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), test, user)
	})

	suite.Run("Unsuccessful Signup", func() {
		suite.Run("Invalid Names", func() {
			test := INVALID_USER_DOMAIN
			test.FirstName = "F"
			test.LastName = "L"
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.SignUp(TEST_CTX, test)
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), test, user)
		})

		suite.Run("Invalid Username", func() {
			test := INVALID_USER_DOMAIN
			test.Credential.Username = "R"
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.SignUp(TEST_CTX, test)
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), test, user)
		})
		suite.Run("Duplicate Username", func() {
			test := INVALID_USER_DOMAIN
			// Assumes that there are duplicate record error
			suite.repo.On("ValidateRepoState", TEST_CTX, test).Return(NewConflictError("User", errors.New("username taken"))).Once()
			// Start sign up test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.SignUp(TEST_CTX, test)
			assert.Error(suite.T(), err, NewConflictError("User", fmt.Errorf("username taken")))
			assert.ErrorContains(suite.T(), err, "conflict state")
			assert.Equal(suite.T(), test, user)
		})
	})
}

func (suite *CredentialUseCaseTestSuite) TestLogin() {
	suite.Run("Successful Login", func() {
		test := VALID_USER_DOMAIN

		// Assumptions for hash password
		suite.service.On("HashPassword", test.Credential.Password).Return("HASHED_PASSWORD").Once()
		// Now test's password is hashed
		test.Credential.Password = "HASHED_PASSWORD"
		// Assumes no error on getting the user by credentials
		suite.repo.On("GetUserByCredentials", TEST_CTX, test.Credential).Return(VALID_USER_DOMAIN_WITH_ID, nil).Once()
		// Now test should have an ID
		test = VALID_USER_DOMAIN_WITH_ID
		// Assumes generating refersh token id flawlessly
		suite.repo.On("SetNewUserToken", TEST_CTX, VALID_USER_DOMAIN_WITH_ID).Return(VALID_USER_TOKENS_ONLY_ID, nil).Once()
		// Now test tokens have an ID
		test.Credential.Tokens = VALID_USER_TOKENS_ONLY_ID
		// Assumes generation run well
		suite.service.On("GenerateUserTokens", test).Return(VALID_USER_TOKENS, nil).Once()
		// Now test has complete tokens
		test.Credential.Tokens = VALID_USER_TOKENS

		// Start login test
		uc := NewCredentialUseCase(suite.repo, suite.service)
		user, err := uc.Login(TEST_CTX, vo.UserCredential{Username: test.Credential.Username, Password: test.Credential.Password})
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), test, user)
	})

	suite.Run("Unsuccessful Login", func() {
		suite.Run("User's Credential Is Not Valid", func() {
			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, vo.UserCredential{Username: "H", Password: "D"})
			assert.Error(suite.T(), err)
			assert.Empty(suite.T(), user)
		})

		suite.Run("User's Credential Not Found", func() {
			test := INVALID_USER_DOMAIN

			// Assumptions for hash password
			suite.service.On("HashPassword", test.Credential.Password).Return("HASHED_PASSWORD").Once()
			// Now test's password is hashed
			test.Credential.Password = "HASHED_PASSWORD"
			// Assumes user's credentials not found in repo
			suite.repo.On("GetUserByCredentials", TEST_CTX, test.Credential).Return(domain.User{}, fmt.Errorf("user not found")).Once()

			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, INVALID_USER_DOMAIN.Credential)
			assert.Error(suite.T(), err)
			assert.ErrorIs(suite.T(), err.(AppError).Type, ErrNotFound)
			assert.ErrorContains(suite.T(), err, "record not found")
			assert.Empty(suite.T(), user)
		})

		suite.Run("Failed while creating token family", func() {
			test := VALID_USER_DOMAIN

			// Assumptions for hash password
			suite.service.On("HashPassword", test.Credential.Password).Return("HASHED_PASSWORD").Once()
			// Now test's password is hashed
			test.Credential.Password = "HASHED_PASSWORD"
			// Assumes no error on getting the user by credentials
			suite.repo.On("GetUserByCredentials", TEST_CTX, test.Credential).Return(VALID_USER_DOMAIN_WITH_ID, nil).Once()
			// Now test should have an ID
			test = VALID_USER_DOMAIN_WITH_ID
			// Assumes generating refersh token id failed
			suite.repo.On("SetNewUserToken", TEST_CTX, VALID_USER_DOMAIN_WITH_ID).Return(vo.UserToken{}, fmt.Errorf("cannot set new token")).Once()

			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, VALID_USER_DOMAIN.Credential)
			assert.Error(suite.T(), err)
			assert.ErrorIs(suite.T(), err.(AppError).Type, ErrUnexpected)
			assert.ErrorContains(suite.T(), err, "something unexpected had occured")
			assert.Empty(suite.T(), user)
		})

		suite.Run("Failed on Token Generation", func() {
			test := VALID_USER_DOMAIN

			// Assumptions for hash password
			suite.service.On("HashPassword", test.Credential.Password).Return("HASHED_PASSWORD").Once()
			// Now test's password is hashed
			test.Credential.Password = "HASHED_PASSWORD"
			// Assumes no error on getting the user by credentials
			suite.repo.On("GetUserByCredentials", TEST_CTX, test.Credential).Return(VALID_USER_DOMAIN_WITH_ID, nil).Once()
			// Now test should have an ID
			test = VALID_USER_DOMAIN_WITH_ID
			// Assumes generating refersh token id flawlessly
			suite.repo.On("SetNewUserToken", TEST_CTX, VALID_USER_DOMAIN_WITH_ID).Return(VALID_USER_TOKENS_ONLY_ID, nil).Once()
			// Now test tokens have an ID
			test.Credential.Tokens = VALID_USER_TOKENS_ONLY_ID
			// Assumes token generation and undoing fails
			suite.service.On("GenerateUserTokens", test).Return(vo.UserToken{}, fmt.Errorf("unable to generate token")).Once()
			suite.repo.On("UndoSetUserToken", TEST_CTX, test.Credential.Tokens.TokenID).Return(fmt.Errorf("unable to delete token id"))

			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, VALID_USER_DOMAIN.Credential)
			assert.Error(suite.T(), err)
			assert.ErrorIs(suite.T(), err.(AppError).Type, ErrUnexpected)
			assert.ErrorContains(suite.T(), err, "something unexpected had occured")
			assert.Empty(suite.T(), user)
		})
	})
}
