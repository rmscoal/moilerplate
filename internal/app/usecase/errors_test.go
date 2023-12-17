package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := NewError("Test", 400, ErrBadRequest, errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrBadRequest.Error())

	err = NewError("Test", 400, errors.New("test"), errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, errors.New("test").Error())
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("Test", errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrConflictState.Error())

	err = NewConflictError("Test", ErrUnexpected)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrUnexpected.Error())
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("Test", errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrNotFound.Error())

	err = NewNotFoundError("Test", ErrUnexpected)
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrUnexpected.Error())
}

func TestNewRepositoryError(t *testing.T) {
	err := NewRepositoryError("Test", errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrUnexpected.Error())
}

func TestNewServiceError(t *testing.T) {
	err := NewServiceError("Test", errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrUnexpected.Error())
}

func TestNewDomainError(t *testing.T) {
	err := NewDomainError("Test", errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrUnprocessableEntity.Error())
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError(errors.New("test error"))
	assert.Error(t, err)
	assert.ErrorAs(t, err, &AppError{})
	assert.ErrorContains(t, err, ErrUnauthorized.Error())
}
