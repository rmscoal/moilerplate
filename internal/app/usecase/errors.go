package usecase

import (
	"errors"
	"strings"
)

var (
	ErrUnexpected          = errors.New("something unexpected had occurred")
	ErrRequestTimeout      = errors.New("query took too long or client cancelled the request")
	ErrInvalidInput        = errors.New("invalid input syntax")
	ErrUnprocessableEntity = errors.New("unable to process entity")
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("record not found")
	ErrConflictState       = errors.New("conflict state")
	ErrUnauthorized        = errors.New("unauthorized action")
	ErrForbidden           = errors.New("forbidden")
	ErrTooManyRequest      = errors.New("too many request")
)

var ErrNameMapper = map[error]string{
	ErrUnexpected:          "UnexpectedError",
	ErrRequestTimeout:      "RequestTimeOutError",
	ErrInvalidInput:        "InvalidInputError",
	ErrUnprocessableEntity: "DomainValidationError",
	ErrBadRequest:          "BadRequestError",
	ErrNotFound:            "NotFoundError",
	ErrConflictState:       "ConflictDuplicationError",
	ErrUnauthorized:        "UnauthorizedError",
	ErrForbidden:           "ForbiddenError",
}

type AppError struct {
	Code    int              `json:"code,omitempty" example:"400"`
	Type    error            `json:"-"`
	Message string           `json:"message,omitempty" example:"this is error message"`
	Errors  []AppErrorDetail `json:"errors,omitempty"`
}

type AppErrorDetail struct {
	Domain  string `json:"domain,omitempty" example:"domain error"`
	Reason  string `json:"reason,omitempty" example:"this is descriptive error reason"`
	Message string `json:"message,omitempty" example:"this is descriptive error message"`
	Report  string `json:"report,omitempty" example:"Please report incident to https://your-report.com"`
}

func (err AppError) Error() string {
	return err.Message
}

func NewError(domain string, code int, errType, err error) error {
	nErr := AppError{
		Code:    code,
		Type:    errType,
		Message: errType.Error(),
	}
	for _, message := range strings.Split(err.Error(), ";") {
		nErr.Errors = append(nErr.Errors, AppErrorDetail{
			Domain:  domain,
			Reason:  ErrNameMapper[errType],
			Message: strings.Trim(message, " "),
		})
	}

	return nErr
}

func NewErrorWithReport(domain string, code int, errType, err error, report string) error {
	nErr := AppError{
		Code:    code,
		Type:    errType,
		Message: errType.Error(),
	}
	for _, message := range strings.Split(err.Error(), ";") {
		nErr.Errors = append(nErr.Errors, AppErrorDetail{
			Domain:  domain,
			Reason:  ErrNameMapper[errType],
			Message: strings.Trim(message, " "),
			Report:  report,
		})
	}

	return nErr
}

func NewClientError(domain string, err error) error {
	return NewError(domain, 400, ErrBadRequest, err)
}

/* - Repository Group Errors - */

// NewConflictError return either 409 for any conflict state on the database
// like unique index and so on or 500 in case the db connection died or others
func NewConflictError(domain string, err error) error {
	switch {
	case errors.Is(err, ErrUnexpected):
		return NewRepositoryError(domain, err)
	default:
		return NewError(domain+" Repository", 409, ErrConflictState, err)
	}
}

// NewNotFoundError return either 404 for resources that are not found in the
// database or 500 in case the db connection died or others
func NewNotFoundError(domain string, err error) error {
	switch {
	case errors.Is(err, ErrUnexpected):
		return NewRepositoryError(domain, err)
	default:
		return NewError(domain+" Repository", 404, ErrNotFound, err)
	}
}

// NewRepositoryError returns a 500 in case the db connection died or others
func NewRepositoryError(domain string, err error) error {
	switch {
	case errors.Is(err, ErrUnexpected):
		return NewError(domain+" Repository", 500, ErrUnexpected, err)
	case errors.Is(err, ErrNotFound):
		return NewError(domain+" Repository", 404, ErrNotFound, err)
	default:
		return NewError(domain+" Repository", 409, ErrConflictState, err)
	}
}

/* - Domain Group Errors - */

func NewDomainError(domain string, err error) error {
	return NewError(domain, 422, ErrUnprocessableEntity, err)
}

func NewUnauthorizedError(err error) error {
	return NewError("", 401, ErrUnauthorized, err)
}

func NewUnauthorizedErrorWithReport(err error) error {
	return NewErrorWithReport("", 401, ErrUnauthorized, err,
		`This might be due to a stolen token or malformed jti. Please go to www.moilorplate.com/change-password to secure your account`,
	)
}

/* - Service Group Errors - */

// NewRateLimittedError returns 429 when it has too many request
func NewTooManyRequest(err error) error {
	return NewError("Rate Limitted", 429, ErrTooManyRequest, err)
}

// NewServiceError returns 500 when a service/third party fails
func NewServiceError(domain string, err error) error {
	return NewError(domain+" Service", 500, ErrUnexpected, err)
}
