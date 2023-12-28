package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

func (m *Middleware) MetricsMiddleware() gin.HandlerFunc {
	attemptsCounter, _ := m.metrics.Int64UpDownCounter("http.server.request_count",
		metric.WithDescription("Number of Requests"),
		metric.WithUnit("Count"),
	)
	totalDuration, _ := m.metrics.Int64Histogram("http.server.request.duration",
		metric.WithDescription("Time taken by request"),
		metric.WithUnit("Milliseconds"),
	)
	activeRequestCounter, _ := m.metrics.Int64UpDownCounter("http.server.active_requests",
		metric.WithDescription("Number of inflight requests"),
		metric.WithUnit("Count"),
	)
	return func(c *gin.Context) {
		attributes := []attribute.KeyValue{
			semconv.HTTPRoute(c.FullPath()),
			semconv.HTTPMethod(c.Request.Method),
			semconv.HTTPScheme(c.Request.Proto),
		}

		ctx := c.Request.Context()

		activeRequestCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		defer activeRequestCounter.Add(ctx, -1, metric.WithAttributes(attributes...))

		start := time.Now()
		defer func() {
			totalDuration.Record(ctx, int64(time.Since(start).Milliseconds()), metric.WithAttributes(attributes...))
			attemptsCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		}()

		c.Next()
	}
}
