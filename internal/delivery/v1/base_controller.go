package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseControllerV1 struct{}

func (bc BaseControllerV1) jsonErrResponse(c *gin.Context, code int, err any) {
	c.AbortWithStatusJSON(code, gin.H{
		"apiVersion": "1.0",
		"error":      err,
	})
}

func (bc BaseControllerV1) Ok(c *gin.Context, obj ...any) {
	if len(obj) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"apiVersion": "1.0",
			"status":     "OK",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"apiVersion": "1.0",
		"status":     "OK",
		"data":       obj[0],
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
	c.JSON(http.StatusOK, gin.H{
		"apiVersion": "1.0",
		"status":     "OK",
		"data":       obj[0],
		"paging":     obj[1],
	})
}

func (bc BaseControllerV1) Created(c *gin.Context, obj ...any) {
	if len(obj) == 0 {
		c.JSON(http.StatusCreated, gin.H{
			"apiVersion": "1.0",
			"status":     "OK",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"apiVersion": "1.0",
		"status":     "OK",
		"data":       obj[0],
	})
}

func (bc BaseControllerV1) ClientError(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusBadRequest, err)
}

func (bc BaseControllerV1) Unauthorized(c *gin.Context) {
	bc.jsonErrResponse(c, http.StatusUnauthorized, map[string]any{
		"apiVersion": "1.0",
		"error": map[string]string{
			"code":    "401",
			"message": "Unauthorized",
		},
	})
}

func (bc BaseControllerV1) NotFound(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusNotFound, err)
}

func (bc BaseControllerV1) Conflict(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusConflict, err)
}

func (bc BaseControllerV1) UnprocessableEntity(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusUnprocessableEntity, err)
}

func (bc BaseControllerV1) UnexpectedError(c *gin.Context, err error) {
	bc.jsonErrResponse(c, http.StatusInternalServerError, err)
}
