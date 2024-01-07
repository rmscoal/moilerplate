package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/internal/domain/vo"
	"github.com/rmscoal/moilerplate/testdata/observability"
	mockrepo "github.com/rmscoal/moilerplate/testdata/repo"
	mockservice "github.com/rmscoal/moilerplate/testdata/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

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

	// Set otel to use dummy/testing tracer provider
	otel.SetTracerProvider(observability.NewTestingTracerProvider())
}

func TestCredentialUseCase(t *testing.T) {
	suite.Run(t, new(CredentialUseCaseTestSuite))
}

func (suite *CredentialUseCaseTestSuite) TestSignup() {
	suite.Run("Successful Signup", func() {
		test := VALID_USER_DOMAIN_WITH_ID
		test.Credential.Password = "HASHED_PASSWORD"

		// Assumes HashPassword call returns the hashed password
		suite.service.On("HashPassword", "PASSWORD").Return([]byte("HASHED_PASSWORD"), nil).Once()
		// Assumes that it passes constraint checks while persisting record
		suite.repo.On("CreateNewUser", TEST_CTX, mock.AnythingOfType("domain.User")).Return(test, nil).Once()
		// Assumes preparing generating refresh token family done successfully
		suite.repo.On("SetNewUserToken", TEST_CTX, mock.AnythingOfType("domain.User")).Return(VALID_USER_TOKENS_ONLY_ID, nil).Once()
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

			// Hash the user's password
			suite.service.On("HashPassword", "PASSWORD").Return([]byte("HASHED_PASSWORD"), nil).Once()
			// Now, we say that there are duplicate error while creating the user
			suite.repo.On("CreateNewUser", TEST_CTX, mock.AnythingOfType("domain.User")).Return(test, errors.New("username")).Once()

			// Start sign up test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.SignUp(TEST_CTX, test)
			assert.Error(suite.T(), err, NewConflictError("User", errors.New("username already exists")))
			assert.ErrorContains(suite.T(), err, "conflict state")
			assert.Equal(suite.T(), test, user)
		})

		suite.Run("Failed Password Hashing", func() {
			test := VALID_USER_DOMAIN

			// Assumes HashPassword call returns the hashed password
			suite.service.On("HashPassword", "PASSWORD").Return([]byte(nil), errors.New("failed to hash password")).Once()

			// Start sign up test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.SignUp(TEST_CTX, test)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "something unexpected had occured")
			assert.Equal(suite.T(), user, test)
		})
	})
}

func (suite *CredentialUseCaseTestSuite) TestLogin() {
	suite.Run("Successful Login", func() {
		test := VALID_USER_DOMAIN
		testWithEncodedHashedPassword := VALID_USER_DOMAIN_WITH_ID
		testWithEncodedHashedPassword.Credential.SetEncodedPasswordFromByte([]byte("HASHED_PASSWORD"))

		// Assumes that test's username exists in repository
		suite.repo.On("GetUserByUsername", TEST_CTX, test.Credential.Username).Return(testWithEncodedHashedPassword, nil).Once()
		// Now test is a valid user domain with ID with an encoded hashed password
		test = testWithEncodedHashedPassword
		// Assumes that the incoming password request matches the hash
		suite.service.On("CompareHashAndPassword", mock.Anything, test.Credential.Password, []byte("HASHED_PASSWORD")).Return(true, nil).Once()
		// Assumes generates refresh token flawlessly
		suite.repo.On("SetNewUserToken", TEST_CTX, testWithEncodedHashedPassword).Return(VALID_USER_TOKENS_ONLY_ID, nil).Once()
		// Now test tokens have an ID
		test.Credential.Tokens = VALID_USER_TOKENS_ONLY_ID
		// Assumes generation run well
		suite.service.On("GenerateUserTokens", test).Return(VALID_USER_TOKENS, nil).Once()
		// Now test has complete tokens
		test.Credential.Tokens = VALID_USER_TOKENS
		// goroutines mocks for generateNewHashMixture
		suite.service.On("HashPassword", test.Credential.Password).Return([]byte("HASHED_PASSWORD"), nil)
		suite.repo.On("RotateUserHashPassword", mock.Anything, domain.User{Id: test.Id, Credential: vo.UserCredential{Password: test.Credential.Password}}).Return(nil)

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

		suite.Run("User's Username Not Found", func() {
			test := INVALID_USER_DOMAIN

			// Assumes username is not found
			suite.repo.On("GetUserByUsername", TEST_CTX, test.Credential.Username).Return(domain.User{}, fmt.Errorf("username not found")).Once()

			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, INVALID_USER_DOMAIN.Credential)
			assert.Error(suite.T(), err)
			assert.ErrorIs(suite.T(), err.(AppError).Type, ErrNotFound)
			assert.ErrorContains(suite.T(), err, "record not found")
			assert.Empty(suite.T(), user)
		})

		suite.Run("Failed to Decode Password", func() {
			test := VALID_USER_DOMAIN
			testWithEncodedHashedPassword := VALID_USER_DOMAIN_WITH_ID
			testWithEncodedHashedPassword.Credential.Password = "BAD_ENCODING"

			// Assumes username is found
			suite.repo.On("GetUserByUsername", TEST_CTX, test.Credential.Username).Return(testWithEncodedHashedPassword, nil).Once()

			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, INVALID_USER_DOMAIN.Credential)
			assert.Error(suite.T(), err)
			assert.ErrorIs(suite.T(), err.(AppError).Type, ErrUnprocessableEntity)
			assert.ErrorContains(suite.T(), err, "unable to process entity")
			assert.Equal(suite.T(), testWithEncodedHashedPassword, user)
		})

		suite.Run("Compare Password and Hash Failed", func() {
			test := INVALID_USER_DOMAIN
			testWithEncodedHashedPassword := VALID_USER_DOMAIN_WITH_ID
			testWithEncodedHashedPassword.Credential.SetEncodedPasswordFromByte([]byte("HASHED_PASSWORD"))

			// Assumes username is found
			suite.repo.On("GetUserByUsername", TEST_CTX, test.Credential.Username).Return(testWithEncodedHashedPassword, nil).Once()
			// Assumes that the password is wrong
			suite.service.On("CompareHashAndPassword", mock.Anything, test.Credential.Password, []byte("HASHED_PASSWORD")).Return(false, fmt.Errorf("password does not match")).Once()

			// Start login test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Login(TEST_CTX, INVALID_USER_DOMAIN.Credential)
			assert.Error(suite.T(), err)
			assert.Equal(suite.T(), testWithEncodedHashedPassword, user)
			assert.ErrorIs(suite.T(), err.(AppError).Type, ErrUnauthorized)
			assert.ErrorContains(suite.T(), err, "unauthorized action")
		})

		suite.Run("Failed Creating Token Family", func() {
			test := VALID_USER_DOMAIN
			testWithEncodedHashedPassword := VALID_USER_DOMAIN_WITH_ID
			testWithEncodedHashedPassword.Credential.SetEncodedPasswordFromByte([]byte("HASHED_PASSWORD"))

			// Assumes that the test's username exists in repository
			suite.repo.On("GetUserByUsername", TEST_CTX, test.Credential.Username).Return(testWithEncodedHashedPassword, nil).Once()
			// Now the test is a valid user with id and encoded hashed password
			test = testWithEncodedHashedPassword
			// Assumes that the password matches
			suite.service.On("CompareHashAndPassword", mock.Anything, VALID_USER_DOMAIN.Credential.Password, []byte("HASHED_PASSWORD")).Return(true, nil).Once()
			// Assumes set new token fails
			suite.repo.On("SetNewUserToken", TEST_CTX, test).Return(vo.UserToken{}, ErrUnexpected).Once()

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
			testWithEncodedHashedPassword := VALID_USER_DOMAIN_WITH_ID
			testWithEncodedHashedPassword.Credential.SetEncodedPasswordFromByte([]byte("HASHED_PASSWORD"))

			// Assumes that the test's username exists in repository
			suite.repo.On("GetUserByUsername", TEST_CTX, test.Credential.Username).Return(testWithEncodedHashedPassword, nil).Once()
			// Now the test is a valid user with id and encoded hashed password
			test = testWithEncodedHashedPassword
			// Assumes that the password matches
			suite.service.On("CompareHashAndPassword", mock.Anything, VALID_USER_DOMAIN.Credential.Password, []byte("HASHED_PASSWORD")).Return(true, nil).Once()
			// Assumes set new token fails
			suite.repo.On("SetNewUserToken", TEST_CTX, test).Return(VALID_USER_TOKENS_ONLY_ID, nil).Once()
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

func (suite *CredentialUseCaseTestSuite) TestAuthenticate() {
	suite.Run("Successful Authentication", func() {
		suite.service.On("VerifyAndParseToken", TEST_CTX, "ACCESS_TOKEN").Return("UUID", nil).Once()

		// Start authenticate test
		uc := NewCredentialUseCase(suite.repo, suite.service)
		user, err := uc.Authenticate(TEST_CTX, "ACCESS_TOKEN")
		assert.Nil(suite.T(), err, "There should be no error")
		if assert.NoError(suite.T(), err, "There should be no error") {
			assert.Equal(suite.T(), domain.User{Id: "UUID"}, user)
		}
	})

	suite.Run("Unsuccessful Authentication", func() {
		suite.service.On("VerifyAndParseToken", TEST_CTX, "ACCESS_TOKEN_INVALID").Return("", fmt.Errorf("error on verifying token")).Once()

		// Start authenticate test
		uc := NewCredentialUseCase(suite.repo, suite.service)
		user, err := uc.Authenticate(TEST_CTX, "ACCESS_TOKEN_INVALID")
		assert.Error(suite.T(), err)
		assert.ErrorContains(suite.T(), err, "unauthorized action")
		assert.Empty(suite.T(), user)
	})
}

func (suite *CredentialUseCaseTestSuite) TestRefresh() {
	suite.Run("Successful Refresh", func() {
		test_user := domain.User{
			Id:        "UUID",
			FirstName: "FIRST_NAME",
			LastName:  "LAST_NAME",
			Credential: vo.UserCredential{
				Username: "USER_NAME",
				Tokens:   VALID_USER_TOKENS,
			},
		}

		suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
		suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
		suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(1, nil).Once()
		suite.repo.On("SetNewUserToken", TEST_CTX, test_user).Return(VALID_USER_TOKENS, nil).Once()
		suite.service.On("GenerateUserTokens", test_user).Return(VALID_USER_TOKENS, nil).Once()

		// Start the refresh test
		uc := NewCredentialUseCase(suite.repo, suite.service)
		user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), test_user, user)
	})

	suite.Run("Unsuccesful Refresh", func() {
		user := domain.User{
			Id:        "UUID",
			FirstName: "FIRST_NAME",
			LastName:  "LAST_NAME",
			Credential: vo.UserCredential{
				Username: "USER_NAME",
				Tokens:   VALID_USER_TOKENS,
			},
		}

		suite.Run("Error parsing refresh token", func() {
			// Assumes fail on parsing token
			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("", errors.New("unable to parse token")).Once()

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "unauthorized action")
		})

		suite.Run("Unable to find jti", func() {
			// Remember this return for 4 times
			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			// Assume JTI is not found
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(domain.User{}, errors.New("not found")).Once()

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "unauthorized action")
		})

		suite.Run("Failed Fetching Latest Token Version", func() {
			test_user := user

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			// Assume repo error
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(0, errors.New("error")).Once()

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "something unexpected had occured")
		})

		suite.Run("JTI Issued Previously_Part 1", func() {
			test_user := user
			test_user.Credential.Tokens.Issued = true

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(1, nil).Once()
			suite.repo.On("DeleteUserTokenFamily", TEST_CTX, test_user).Return(nil).Once()

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "unauthorized action")

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(1, nil).Once()
			suite.repo.On("DeleteUserTokenFamily", TEST_CTX, test_user).Return(nil).Once()

			// Start refresh test
			uc = NewCredentialUseCase(suite.repo, suite.service)
			user, err = uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "unauthorized action")
		})

		suite.Run("JTI Issued Previously_Part 2", func() {
			test_user := user

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			// Database version number > user's verion number
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(2, nil).Once()
			suite.repo.On("DeleteUserTokenFamily", TEST_CTX, test_user).Return(nil).Once()

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "unauthorized action")

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			// Database version number > user's verion number
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(2, nil).Once()
			suite.repo.On("DeleteUserTokenFamily", TEST_CTX, test_user).Return(errors.New("")).Once()

			// Start refresh test
			uc = NewCredentialUseCase(suite.repo, suite.service)
			user, err = uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "unauthorized action")
		})

		suite.Run("Failed Setting User Token", func() {
			test_user := user

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(1, nil).Once()
			suite.repo.On("SetNewUserToken", TEST_CTX, test_user).Return(vo.UserToken{}, errors.New("")).Once()

			// Uncomment this if you really need to:
			// suite.repo.AssertExpectations(suite.T())

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "something unexpected had occured")
		})

		suite.Run("Failed Generate Refresh Token", func() {
			test_user := user

			suite.service.On("VerifyAndParseRefreshToken", TEST_CTX, "REFRESH_TOKEN").Return("SOME_JTI", nil).Once()
			suite.repo.On("GetUserByJti", TEST_CTX, "SOME_JTI").Return(test_user, nil).Once()
			suite.repo.On("GetLatestUserTokenVersion", TEST_CTX, test_user).Return(1, nil).Once()
			suite.repo.On("SetNewUserToken", TEST_CTX, test_user).Return(VALID_USER_TOKENS, nil).Once()
			suite.service.On("GenerateUserTokens", test_user).Return(vo.UserToken{}, errors.New("")).Once()
			suite.repo.On("UndoSetUserToken", TEST_CTX, "TOKEN_ID").Return(nil).Once()

			// Start refresh test
			uc := NewCredentialUseCase(suite.repo, suite.service)
			user, err := uc.Refresh(TEST_CTX, "REFRESH_TOKEN")
			assert.Empty(suite.T(), user)
			assert.Error(suite.T(), err)
			assert.ErrorContains(suite.T(), err, "something unexpected had occured")
		})
	})
}
