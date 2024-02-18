package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/moilerplate/internal/app/usecase"
)

type AuthHeader struct {
	BearerToken string `header:"Authorization" binding:"required"`
}

func (m *Middleware) AuthMiddleware(uc usecase.ICredentialUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth AuthHeader
		if err := c.ShouldBindHeader(&auth); err != nil {
			m.Unauthorized(c, usecase.NewUnauthorizedError(errors.New("header not found")))
			return
		}
		auth.BearerToken = strings.ReplaceAll(auth.BearerToken, "Bearer ", "")
		user, err := uc.Authenticate(c.Request.Context(), auth.BearerToken)
		if err != nil {
			m.Unauthorized(c, err)
			return
		}

		m.addToContext(c, "userId", user.ID)
		c.Next()
	}
}
