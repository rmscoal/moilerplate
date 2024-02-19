package repo

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
	"github.com/rmscoal/moilerplate/internal/domain"
	"github.com/rmscoal/moilerplate/testing/observability"
	mockrepo "github.com/rmscoal/moilerplate/testing/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
)

type credentialRepoTestSuite struct {
	suite.Suite
	repo *credentialRepo
	mock sqlmock.Sqlmock
	ctx  context.Context
}

func TestCredentialRepoSuite(t *testing.T) {
	suite.Run(t, new(credentialRepoTestSuite))
}

func (suite *credentialRepoTestSuite) SetupSuite() {
	otel.SetTracerProvider(observability.NewTestingTracerProvider())
	suite.ctx = context.Background()
}

func (suite *credentialRepoTestSuite) SetupTest() {
	// Init sqlmock
	_, gormdb, mock, err := mockrepo.InitGormMock()
	if err != nil {
		suite.T().Fatalf("Error while initializing sqlmock %s", err)
	}

	InitBaseRepo(gormdb, true) // Skipping registry
	suite.repo = NewCredentialRepo()
	suite.mock = mock
	suite.ctx = context.Background()
}

func (suite *credentialRepoTestSuite) TestCreateUser_Success() {
	mock := suite.mock

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("id"))
	mock.ExpectCommit()

	user, err := suite.repo.CreateUser(context.Background(), domain.User{})
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), user) // id and timestamps
}

func (suite *credentialRepoTestSuite) TestCreateUser_Fail_Duplicate() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnError(&pgconn.PgError{Code: DuplicateError.String(), ConstraintName: "some_name"})
	suite.mock.ExpectRollback()

	user, err := suite.repo.CreateUser(context.Background(), domain.User{})
	assert.ErrorContains(suite.T(), err, "already exists")
	assert.NotEmpty(suite.T(), user) // timestamps
}

func (suite *credentialRepoTestSuite) TestCredentialRepo_CreateUser_Fail_ConnDone() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnError(sql.ErrConnDone)
	suite.mock.ExpectRollback()

	user, err := suite.repo.CreateUser(context.Background(), domain.User{})
	assert.ErrorContains(suite.T(), err, usecase.ErrUnexpected.Error())
	assert.NotEmpty(suite.T(), user) // timestamps
}

func (suite *credentialRepoTestSuite) TestCreateUser_Fail_Unknown() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnError(errors.New("unknown error"))
	suite.mock.ExpectRollback()

	user, err := suite.repo.CreateUser(context.Background(), domain.User{})
	assert.ErrorContains(suite.T(), err, usecase.ErrUnexpected.Error())
	assert.NotEmpty(suite.T(), user) // timestamps
}

func (suite *credentialRepoTestSuite) TestGetUserByUsername_Success() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE username (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow("id", "username"))

	user, err := suite.repo.GetUserByUsername(suite.ctx, "username")
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), user)
}

func (suite *credentialRepoTestSuite) TestGetUserByUsername_Fail_ConnDone() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE username (.+)`).WillReturnError(sql.ErrConnDone)

	user, err := suite.repo.GetUserByUsername(suite.ctx, "username")
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), user)
}

func (suite *credentialRepoTestSuite) TestGetUserByUsername_Fail_NotFound() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE username (.+)`).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))

	user, err := suite.repo.GetUserByUsername(suite.ctx, "username")
	assert.ErrorContains(suite.T(), err, "not found")
	assert.Empty(suite.T(), user)
}

func (suite *credentialRepoTestSuite) TestGetUserByID_Success() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE id (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).AddRow("id", "username"))

	user, err := suite.repo.GetUserByID(suite.ctx, "username")
	assert.Nil(suite.T(), err)
	assert.NotEmpty(suite.T(), user)
}

func (suite *credentialRepoTestSuite) TestGetUserByID_Fail_ConnDone() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE id (.+)`).WillReturnError(sql.ErrConnDone)

	user, err := suite.repo.GetUserByID(suite.ctx, "username")
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), user)
}

func (suite *credentialRepoTestSuite) TestGetUserByID_Fail_NotFound() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE id (.+)`).WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))

	user, err := suite.repo.GetUserByID(suite.ctx, "username")
	assert.ErrorContains(suite.T(), err, "not found")
	assert.Empty(suite.T(), user)
}
