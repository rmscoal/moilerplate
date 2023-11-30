package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
)

type BaseControllerV1 struct{}

type Data struct {
	ApiVersion string `json:"apiVersion,omitempty" example:"1.0"`
	Status     string `json:"status,omitempty" example:"OK"`
	Data       any    `json:"data,omitempty"`
	Paging     any    `json:"paging,omitempty" extensions:"x-nullable,x-omitempty"`
}

type Error struct {
	ApiVersion string `json:"apiVersion,omitempty" example:"1.0"`
	Error      any    `json:"error,omitempty"`
}

func (bc BaseControllerV1) jsonErrResponse(c *gin.Context, code int, err any) {
	c.AbortWithStatusJSON(code, Error{
		ApiVersion: "1.0",
		Error:      err,
	})
}

func (bc BaseControllerV1) Ok(c *gin.Context, obj ...any) {
	if len(obj) == 0 {
		c.JSON(http.StatusOK, Data{
			ApiVersion: "1.0",
			Status:     "OK",
		})
		return
	}
	c.JSON(http.StatusOK, Data{
		ApiVersion: "1.0",
		Status:     "OK",
		Data:       obj[0],
	})
}

// OkWithPage sends a response with a 200 http status code
// to the client. It sends the data with a paging information.
// It requires two data to be sent. If there are less than
// two data, it calls `Ok` internally.
func (bc BaseControllerV1) OkWithPage(c *gin.Context, obj ...any) {
	if len(obj) < 2 {
		bc.Ok(c, obj)
		return
	}
	c.JSON(http.StatusOK, Data{
		ApiVersion: "1.0",
		Status:     "OK",
		Data:       obj[0],
		Paging:     obj[1],
	})
}

func (bc BaseControllerV1) Created(c *gin.Context, obj ...any) {
	if len(obj) == 0 {
		c.JSON(http.StatusCreated, Data{
			ApiVersion: "1.0",
			Status:     "OK",
		})
		return
	}
	c.JSON(http.StatusCreated, Data{
		ApiVersion: "1.0",
		Status:     "OK",
		Data:       obj[0],
	})
}

func (bc BaseControllerV1) ClientError(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusBadRequest, err)
}

func (bc BaseControllerV1) Unauthorized(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusUnauthorized, err)
}

func (bc BaseControllerV1) NotFound(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusNotFound, err)
}

func (bc BaseControllerV1) RequestTimeout(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusRequestTimeout, err)
}

func (bc BaseControllerV1) Conflict(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusConflict, err)
}

func (bc BaseControllerV1) UnprocessableEntity(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusUnprocessableEntity, err)
}

func (bc BaseControllerV1) TooManyRequest(c *gin.Context) {
	bc.jsonErrResponse(c, http.StatusTooManyRequests, gin.H{"message": "Too many request, try again later"})
}

func (bc BaseControllerV1) UnexpectedError(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusInternalServerError, err)
}

func (bc BaseControllerV1) SummariesUseCaseError(c *gin.Context, err any) {
	if appError, ok := err.(usecase.AppError); ok {
		switch appError.Type {
		case usecase.ErrUnexpected:
			bc.UnexpectedError(c, appError)
		case usecase.ErrRequestTimeout:
			bc.RequestTimeout(c, appError)
		case usecase.ErrInvalidInput:
			bc.ClientError(c, appError)
		case usecase.ErrUnprocessableEntity:
			bc.UnprocessableEntity(c, appError)
		case usecase.ErrBadRequest:
			bc.ClientError(c, appError)
		case usecase.ErrNotFound:
			bc.NotFound(c, appError)
		case usecase.ErrConflictState:
			bc.Conflict(c, appError)
		case usecase.ErrUnauthorized:
			bc.Unauthorized(c, appError)
		}
	} else {
		bc.UnexpectedError(c, usecase.AppError{
			Code:    500,
			Message: "unable to display message",
			Errors: []usecase.AppErrorDetail{
				{
					Message: "unable to display error message",
					Report:  "Please contact admin@example.com regarding this issue",
				},
			},
		})
	}
}
