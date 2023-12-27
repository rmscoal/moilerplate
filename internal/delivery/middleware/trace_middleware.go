package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func (m *Middleware) TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := m.tracer.Start(c.Request.Context(), c.FullPath(), trace.WithAttributes(
			attribute.String(string(semconv.HTTPRequestMethodKey), c.Request.Method),
			attribute.String(string(semconv.HTTPRouteKey), c.FullPath()),
			attribute.String(string(semconv.HTTPSchemeKey), c.Request.Proto),
		))
		defer span.End()

		// Hot reload please work
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
