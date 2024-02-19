package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rmscoal/moilerplate/pkg/doorkeeper"
	"github.com/rmscoal/moilerplate/testing/observability"
	mockrepo "github.com/rmscoal/moilerplate/testing/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type DoorkeeperServiceTestSuite struct {
	suite.Suite
	service *doorkeeperService
	dk      *doorkeeper.Doorkeeper
	db      *sql.DB
	mock    sqlmock.Sqlmock
	ctx     context.Context
}

func TestDoorkeeperService(t *testing.T) {
	suite.Run(t, new(DoorkeeperServiceTestSuite))
}

func (suite *DoorkeeperServiceTestSuite) SetupSuite() {
	otel.SetTracerProvider(observability.NewTestingTracerProvider())

	dk := doorkeeper.GetDoorkeeper(
		// JWT
		doorkeeper.RegisterJWTIssuer("TESTING"),
		doorkeeper.RegisterJWTSignMethod("HMAC", "256"),
		doorkeeper.RegisterJWTPublicKey("verysecretkey"),
		doorkeeper.RegisterJWTPrivateKey("verysecretkey"),
		doorkeeper.RegisterJWTAccessDuration(5*time.Minute),
		doorkeeper.RegisterJWTRefreshDuration(10*time.Minute),
		// Encryptor
		doorkeeper.RegisterEncryptorSecretKey("verystrongsecretkey"),
		// General
		doorkeeper.RegisterGeneralHasherFunc("SHA384"),
	)

	db, _, mock, err := mockrepo.InitGormMock()
	if err != nil {
		suite.T().Fatalf("failed to initialize mock repo")
	}

	suite.dk = dk
	suite.db = db
	suite.mock = mock
	suite.ctx = context.Background()
	suite.service = NewDoorkeeperService(suite.dk, suite.db)
}

func (suite *DoorkeeperServiceTestSuite) SetupTest() {}

func (suite *DoorkeeperServiceTestSuite) TestHashAndEncodeStringWithSalt_Success() {
	result := suite.service.HashAndEncodeStringWithSalt(context.Background(), "password", "salt")
	assert.NotEmpty(suite.T(), result)
}

func (suite *DoorkeeperServiceTestSuite) TestComparePasswords_Success() {
	match, err := suite.service.ComparePasswords(suite.ctx, suite.service.HashAndEncodeStringWithSalt(suite.ctx, "password", "salt"), "password", "salt")
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), match)
}

func (suite *DoorkeeperServiceTestSuite) TestComparePasswords_Fail_Decode() {
	match, err := suite.service.ComparePasswords(suite.ctx, "wrong padding", "password", "salt")
	assert.Error(suite.T(), err)
	assert.False(suite.T(), match)
}

func (suite *DoorkeeperServiceTestSuite) TestComparePasswords_Fail_Mismatch() {
	match, err := suite.service.ComparePasswords(suite.ctx, suite.service.HashAndEncodeStringWithSalt(suite.ctx, "password", "salt"), "wrong", "salt")
	assert.Nil(suite.T(), err)
	assert.False(suite.T(), match)
}

func (suite *DoorkeeperServiceTestSuite) TestGenerateTokens_Success() {
	suite.mock.ExpectExec("INSERT INTO access_versionings (.+) VALUES (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(sqlmock.AnyArg(), "prevJTI", "subject")

	token, err := suite.service.GenerateTokens(suite.ctx, "subject", nil)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *DoorkeeperServiceTestSuite) TestGenerateTokens_Fail_Query() {
	suite.mock.ExpectExec("INSERT INTO access_versionings (.+) VALUES (.+)").WillReturnError(sql.ErrConnDone)

	prevJTI := new(string)
	*prevJTI = "prevJTI"

	token, err := suite.service.GenerateTokens(suite.ctx, "subject", prevJTI)
	assert.ErrorContains(suite.T(), err, sql.ErrConnDone.Error())
	assert.Empty(suite.T(), token)
}

func (suite *DoorkeeperServiceTestSuite) TestValidateAccessToken_Success() {
	// Mock for generate access token
	suite.mock.ExpectExec("INSERT INTO access_versionings (.+) VALUES (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(sqlmock.AnyArg(), "prevJTI", "subject")
	prevJTI := new(string)
	*prevJTI = "prevJTI"

	// Generate access token
	token, _ := suite.service.GenerateTokens(suite.ctx, "subject", prevJTI)

	// Test
	userID, err := suite.service.ValidateAccessToken(suite.ctx, token.AccessToken)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "subject", userID)
}

func (suite *DoorkeeperServiceTestSuite) TestValidateAccessToken_Fail_Parse() {
	userID, err := suite.service.ValidateAccessToken(suite.ctx, "some_random_string")
	assert.ErrorContains(suite.T(), err, "token is malformed")
	assert.Empty(suite.T(), userID)
}

func (suite *DoorkeeperServiceTestSuite) TestValidateRefreshToken_Success() {
	// Mock for generate access token
	suite.mock.ExpectExec("INSERT INTO access_versionings (.+) VALUES (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(sqlmock.AnyArg(), nil, "user_id")

	// Mock for test
	suite.mock.ExpectQuery("SELECT av1.user_id FROM access_versionings av1 WHERE (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow("user_id"))
	suite.mock.ExpectExec("INSERT INTO access_versionings (.+) VALUES (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "user_id")

	// Generate refresh token
	token, _ := suite.service.GenerateTokens(suite.ctx, "user_id", nil)

	// Test
	token, err := suite.service.ValidateRefreshToken(suite.ctx, token.RefreshToken)
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *DoorkeeperServiceTestSuite) TestValidateRefreshToken_Fail_Parse() {
	// Test
	token, err := suite.service.ValidateRefreshToken(suite.ctx, "some_random_string")
	assert.ErrorContains(suite.T(), err, "token is malformed")
	assert.Empty(suite.T(), token)
}

func (suite *DoorkeeperServiceTestSuite) TestValidateRefreshToken_Fail_QueryRow() {
	// Mock for generate access token
	suite.mock.ExpectExec("INSERT INTO access_versionings (.+) VALUES (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1)).WithArgs(sqlmock.AnyArg(), nil, "user_id")

	suite.mock.ExpectQuery("SELECT av1.user_id FROM access_versionings av1 WHERE (.+)").WillReturnError(sql.ErrNoRows)

	// Generate refresh token
	token, _ := suite.service.GenerateTokens(suite.ctx, "user_id", nil)

	// Test
	token, err := suite.service.ValidateRefreshToken(suite.ctx, token.RefreshToken)
	assert.ErrorContains(suite.T(), err, ErrTokenExpiredOrInvalidated.Error())
	assert.Empty(suite.T(), token)
}
