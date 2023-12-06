package middleware

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/rmscoal/moilerplate/internal/app/service"
)

// Thanks to: https://blog.logrocket.com/rate-limiting-go-application/
func (m *Middleware) RateLimiterMiddleware(service service.IRaterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			m.UnexpectedError(c, fmt.Errorf("unable to identify client ip"))
			return
		}
		if service.IsClientAllowed(c.Request.Context(), ip) {
			c.Next()
		} else {
			m.TooManyRequest(c)
		}
	}
}
