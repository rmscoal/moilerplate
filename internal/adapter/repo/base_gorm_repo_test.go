package repo

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
	mockrepo "github.com/rmscoal/moilerplate/testing/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var (
	testErrForeignKey = &pgconn.PgError{
		Code: "23503",
	}

	testErrDuplicated = &pgconn.PgError{
		Code: "23505",
	}
)

type baseRepoImplTestSuite struct {
	suite.Suite
	sqldb *sql.DB
	repo  *baseRepo
	mock  sqlmock.Sqlmock
}

func (suite *baseRepoImplTestSuite) SetupTest() {
	// Init sqlmock
	sqldb, gormdb, mock, err := mockrepo.InitGormMock()
	if err != nil {
		suite.T().Fatalf("Error while initializing sqlmock %s", err)
	}

	InitBaseRepo(gormdb, false) // Skipping registry
	suite.sqldb = sqldb
	suite.mock = mock
	suite.repo = gormRepo
}

func (suite *baseRepoImplTestSuite) TearDownTest() {
	suite.sqldb.Close()
}

// ------ TESTING SECTION ------
func TestBaseRepoImpl(t *testing.T) {
	suite.Run(t, new(baseRepoImplTestSuite))
}

func (suite *baseRepoImplTestSuite) TestRegisterIndexes_Success() {
	mock := suite.mock

	mock.ExpectQuery("SELECT (.+) FROM pg_indexes JOIN (.+) JOIN (.+) WHERE (.+) GROUP BY (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"index_name", "indexed_columns"}).AddRow("yeah", "oh"))

	err := suite.repo.registerIndexes()
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), mock.ExpectationsWereMet())
}

func (suite *baseRepoImplTestSuite) TestRegisterIndexes_Fail_ConnClosed() {
	mock := suite.mock

	mock.ExpectQuery("SELECT (.+) FROM pg_indexes JOIN (.+) JOIN (.+) WHERE (.+) GROUP BY (.+)").WillReturnError(sql.ErrConnDone)

	err := suite.repo.registerIndexes()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), mock.ExpectationsWereMet())
}

func (suite *baseRepoImplTestSuite) TestRegisterForeignKeys_Success() {
	mock := suite.mock

	mock.ExpectQuery("SELECT (.+) FROM pg_constraint JOIN (.+) JOIN (.+) WHERE (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"foreign_key_name", "referenced_table"}).AddRow("yeah", "oh"))

	err := suite.repo.registerForeignKeys()
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), mock.ExpectationsWereMet())
}

func (suite *baseRepoImplTestSuite) TestRegisterForeignKeys_Fail_ConnClosed() {
	mock := suite.mock

	mock.ExpectQuery("SELECT (.+) FROM pg_constraint JOIN (.+) JOIN (.+) WHERE (.+)").WillReturnError(sql.ErrConnDone)

	err := suite.repo.registerForeignKeys()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), mock.ExpectationsWereMet())
}

func (suite *baseRepoImplTestSuite) TestDetectConstraintError_NoError() {
	assert.Nil(suite.T(), suite.repo.DetectConstraintError(nil))
}

func (suite *baseRepoImplTestSuite) TestDetectConstraintError_NotPgError() {
	err := suite.repo.DetectConstraintError(errors.New("hello"))
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, usecase.ErrUnexpected.Error())
}

func (suite *baseRepoImplTestSuite) TestDetectConstraintError_PgError_DuplicateError() {
	err := suite.repo.DetectConstraintError(testErrDuplicated)
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "already exists")
}

func (suite *baseRepoImplTestSuite) TestDetectConstraintError_PgError_ForeignKeyError() {
	err := suite.repo.DetectConstraintError(testErrForeignKey)
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, "association error to")
}

func (suite *baseRepoImplTestSuite) TestDetectNotFoundError_NoError() {
	err := suite.repo.DetectNotFoundError(nil)
	assert.Nil(suite.T(), err)
}

func (suite *baseRepoImplTestSuite) TestDetectNotFoundError_Unknown() {
	err := suite.repo.DetectNotFoundError(errors.New("random error"))
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, usecase.ErrUnexpected.Error())
}

func (suite *baseRepoImplTestSuite) TestDetectNotFoundError_NotFound() {
	err := suite.repo.DetectNotFoundError(gorm.ErrRecordNotFound)
	assert.Error(suite.T(), err)
	assert.ErrorContains(suite.T(), err, usecase.ErrNotFound.Error())
}
