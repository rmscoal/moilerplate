package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/go-restful-monolith-boilerplate/internal/app/usecase"
)

type AuthHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func (m *Middleware) AuthMiddleware(uc usecase.ICredentialUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth AuthHeader
		if err := c.ShouldBindHeader(&auth); err != nil {
			m.Unauthorized(c, usecase.NewUnauthorizedError(fmt.Errorf("header not found")))
			return
		}
		auth.Authorization = strings.ReplaceAll(auth.Authorization, "Bearer ", "")
		user, err := uc.Authorize(c.Request.Context(), auth.Authorization)
		if err != nil {
			m.Unauthorized(c, err)
			return
		}

		m.addToContext(c, "userId", user.Id)
		c.Next()
	}
}

func (m *Middleware) AdminMiddleware(uc usecase.ICredentialUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, err := c.Cookie("x-admin-key")
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "auth")
			return
		}

		if err := uc.AuthAdmin(c.Request.Context(), admin); err != nil {
			m.Unauthorized(c, err)
			return
		}

		c.Next()
	}
}
