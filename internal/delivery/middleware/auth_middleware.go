package middleware

import (
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
		if err := c.Copy().ShouldBindHeader(&auth); err != nil {
			m.Unauthorized(c, err)
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
