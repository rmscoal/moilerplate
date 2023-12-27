package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (m *Middleware) TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := m.tracer.Start(c.Request.Context(), c.FullPath(), trace.WithAttributes(
			attribute.String("http-method", c.Request.Method),
		))
		defer span.End()

		// Hot reload please work
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
