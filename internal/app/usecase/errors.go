package usecase

import (
	"errors"
	"strings"
)

var (
	ErrUnexpected          = errors.New("something unexpected had occured")
	ErrRequestTimeout      = errors.New("query took too long or client cancelled the request")
	ErrInvalidInput        = errors.New("invalid input syntax")
	ErrUnprocessableEntity = errors.New("unable to process entity")
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("record not found")
	ErrConflictState       = errors.New("conflict state")
	ErrUnauthorized        = errors.New("unauthorized action")
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

func NewConflictError(domain string, err error) error {
	return NewError(domain, 409, ErrConflictState, err)
}

func NewDomainError(domain string, err error) error {
	return NewError(domain, 422, ErrUnprocessableEntity, err)
}

func NewRepositoryError(domain string, err error) error {
	domain = domain + " Repository"
	switch {
	case errors.Is(err, ErrUnexpected):
		return NewError(domain, 500, ErrUnexpected, err)
	case errors.Is(err, ErrNotFound):
		return NewNotFoundError(domain, err)
	default:
		return NewConflictError(domain, err)
	}
}

func NewServiceError(domain string, err error) error {
	return NewError(domain+" Service", 500, ErrUnexpected, err)
}

func NewNotFoundError(domain string, err error) error {
	return NewError(domain, 404, ErrNotFound, err)
}

func NewUnauthorizedError(err error) error {
	return NewError("", 401, ErrUnauthorized, err)
}

func NewUnauthorizedErrorWithReport(err error) error {
	return NewErrorWithReport("", 401, ErrUnauthorized, err,
		`This might be due to a stolen token or malformed jti. Please go to www.moilorplate.com/change-password to secure your account`,
	)
}
