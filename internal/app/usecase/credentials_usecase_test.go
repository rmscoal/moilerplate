package usecase

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain"
// 	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/domain/vo"
// 	mockrepo "github.com/rmscoal/go-restful-monolith-boilerplate/testdata/repo"
// 	mockservice "github.com/rmscoal/go-restful-monolith-boilerplate/testdata/service"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/suite"
// )

// /*
// Scenario:

// A new user registers via its,
// 1. firstname,
// 2. lastname,
// 3. username,
// 4. password,
// 5. email, and
// 6. phone number

// The creation of a new user is valid if:
// a. the firstname is not empty and is around 3 - 20 bytes long
// b. the lastname is not empty and is around 3 - 25 bytes long
// c. the username is not empty and is around 3 - 20 bytes long
// d. the email is not empty and is a valid email checked by regex
// e. the password is not empty
// f. the phone number is not empty and is a valid indonesian phone number
// */

// var (
// 	VALID_USER_DOMAIN = domain.User{
// 		FirstName: "RMSCOAL",
// 		LastName:  "RMSCOAL",
// 		Emails: []vo.UserEmail{
// 			{
// 				Email:     "rmscoal@test.com",
// 				IsPrimary: true,
// 			},
// 		},
// 		PhoneNumber: "0812345699",
// 		Credential: vo.UserCredential{
// 			Username: "RMSCOAL",
// 			Password: "PASSWORD",
// 		},
// 	}
// 	VALID_USER_DOMAIN_WITH_ID = domain.User{
// 		Id:        "UNIQUE_ID",
// 		FirstName: "RMSCOAL",
// 		LastName:  "RMSCOAL",
// 		Emails: []vo.UserEmail{
// 			{
// 				Email:     "rmscoal@test.com",
// 				IsPrimary: true,
// 			},
// 		},
// 		PhoneNumber: "0812345699",
// 		Credential: vo.UserCredential{
// 			Username: "RMSCOAL",
// 			Password: "PASSWORD",
// 		},
// 	}
// 	// Use this variable for testing unsuccessful cases.
// 	// It makes readibility better and understanable.
// 	INVALID_USER_DOMAIN = VALID_USER_DOMAIN
// )

// type CredentialUseCaseTestSuite struct {
// 	suite.Suite
// 	repo    *mockrepo.CredentialRepoMock
// 	service *mockservice.DoorkeeperServiceMock
// }

// func (suite *CredentialUseCaseTestSuite) SetupTest() {
// 	suite.repo = new(mockrepo.CredentialRepoMock)
// 	suite.service = new(mockservice.DoorkeeperServiceMock)
// }

// func TestCredentialUseCase(t *testing.T) {
// 	suite.Run(t, new(CredentialUseCaseTestSuite))
// }

// func (suite *CredentialUseCaseTestSuite) TestSignup() {
// 	// Setup context for test
// 	ctx := context.Background()

// 	suite.Run("Successful Signup", func() {
// 		test := VALID_USER_DOMAIN_WITH_ID
// 		test.Credential.Password = "HASHED_PASSWORD"

// 		// Assumes that it passes all validity repo state checks
// 		suite.repo.On("ValidateRepoState", ctx, VALID_USER_DOMAIN).Return(nil).Once()
// 		// Assumes that it passes constraint checks while persisting record
// 		suite.repo.On("CreateNewUser", ctx, mock.AnythingOfType("domain.User")).Return(test, nil).Once()
// 		// Assumes HashPassword call returns the hashed password
// 		suite.service.On("HashPassword", mock.AnythingOfType("string")).Return("HASHED_PASSWORD").Once()
// 		// Assumes HashPassword call successfully return the hashed password
// 		suite.service.On("GenerateToken", mock.AnythingOfType("domain.User")).Return("TOKEN", nil).Once()

// 		test.Credential.Token = "TOKEN"

// 		// Start sign up test
// 		uc := NewCredentialUseCase(suite.repo, suite.service)
// 		user, err := uc.SignUp(ctx, VALID_USER_DOMAIN)
// 		assert.Nil(suite.T(), err)
// 		assert.Equal(suite.T(), test, user)
// 	})

// 	suite.Run("Unsuccessful Signup", func() {
// 		suite.Run("Invalid Names", func() {
// 			test := INVALID_USER_DOMAIN
// 			test.FirstName = "F"
// 			test.LastName = "L"
// 			uc := NewCredentialUseCase(suite.repo, suite.service)
// 			user, err := uc.SignUp(ctx, test)
// 			assert.Error(suite.T(), err)
// 			assert.Equal(suite.T(), test, user)
// 		})

// 		suite.Run("Invalid Username", func() {
// 			test := INVALID_USER_DOMAIN
// 			test.Credential.Username = "R"
// 			uc := NewCredentialUseCase(suite.repo, suite.service)
// 			user, err := uc.SignUp(ctx, test)
// 			assert.Error(suite.T(), err)
// 			assert.Equal(suite.T(), test, user)
// 		})
// 		suite.Run("Duplicate Username", func() {
// 			test := INVALID_USER_DOMAIN
// 			// Assumes that there are duplicate record error
// 			suite.repo.On("ValidateRepoState", ctx, test).Return(NewConflictError("User", errors.New("username taken"))).Once()
// 			// Start sign up test
// 			uc := NewCredentialUseCase(suite.repo, suite.service)
// 			user, err := uc.SignUp(ctx, test)
// 			assert.Error(suite.T(), err, NewConflictError("User", fmt.Errorf("username taken")))
// 			assert.ErrorContains(suite.T(), err, "conflict state")
// 			assert.Equal(suite.T(), test, user)
// 		})
// 	})
// }
