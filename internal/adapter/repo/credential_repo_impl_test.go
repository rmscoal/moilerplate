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
	mockrepo "github.com/rmscoal/moilerplate/testing/repo"
	"github.com/stretchr/testify/assert"
)

func TestCredentialRepo_CreateUser_Success(t *testing.T) {
	_, gorm, mock, err := mockrepo.InitGormMock()
	if err != nil {
		t.Fatalf("failed to initialize mock repo")
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("id"))
	mock.ExpectCommit()

	InitBaseRepo(gorm, false)
	repo := NewCredentialRepo()
	user, err := repo.CreateUser(context.Background(), domain.User{})
	assert.Nil(t, err)
	assert.NotEmpty(t, user) // id and timestamps
}

func TestCredentialRepo_CreateUser_Fail_Duplicate(t *testing.T) {
	_, gorm, mock, err := mockrepo.InitGormMock()
	if err != nil {
		t.Fatalf("failed to initialize mock repo")
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnError(&pgconn.PgError{Code: DuplicateError.String(), ConstraintName: "some_name"})
	mock.ExpectRollback()

	InitBaseRepo(gorm, false)
	repo := NewCredentialRepo()
	user, err := repo.CreateUser(context.Background(), domain.User{})
	assert.ErrorContains(t, err, "already exists")
	assert.NotEmpty(t, user) // timestamps
}

func TestCredentialRepo_CreateUser_Fail_ConnDone(t *testing.T) {
	_, gorm, mock, err := mockrepo.InitGormMock()
	if err != nil {
		t.Fatalf("failed to initialize mock repo")
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	InitBaseRepo(gorm, false)
	repo := NewCredentialRepo()
	user, err := repo.CreateUser(context.Background(), domain.User{})
	assert.ErrorContains(t, err, usecase.ErrUnexpected.Error())
	assert.NotEmpty(t, user) // timestamps
}

func TestCredentialRepo_CreateUser_Fail_Unknown(t *testing.T) {
	_, gorm, mock, err := mockrepo.InitGormMock()
	if err != nil {
		t.Fatalf("failed to initialize mock repo")
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+) RETURNING (.+)`).
		WillReturnError(errors.New("unknown error"))
	mock.ExpectRollback()

	InitBaseRepo(gorm, false)
	repo := NewCredentialRepo()
	user, err := repo.CreateUser(context.Background(), domain.User{})
	assert.ErrorContains(t, err, usecase.ErrUnexpected.Error())
	assert.NotEmpty(t, user) // timestamps
}
