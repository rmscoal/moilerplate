package repo

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGormMock with the help of sqlmock is intended to create a new gorm.DB
// mock for unit testing purposes.
func InitGormMock() (*sql.DB, *gorm.DB, sqlmock.Sqlmock, error) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return sqldb, gormdb, mock, nil
}

type BaseRepoMock struct {
	mock.Mock
}

func (repo *BaseRepoMock) DetectConstraintError(err error) error {
	args := repo.Called(err)
	return args.Error(0)
}

func (repo *BaseRepoMock) DetectNotFoundError(err error) error {
	args := repo.Called(err)
	return args.Error(0)
}
