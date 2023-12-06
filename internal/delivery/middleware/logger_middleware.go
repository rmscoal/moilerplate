package middleware

import (
	"time"

	"github.com/rmscoal/moilerplate/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RequestLog struct {
	Latency      time.Duration
	StatusCode   int
	ClientIP     string
	Method       string
	RelativePath string
	UserAgent    string
}

func LogRequestMiddleware(logger *logger.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)
		requestLog := RequestLog{
			Latency:      latency,
			StatusCode:   c.Writer.Status(),
			ClientIP:     c.ClientIP(),
			Method:       c.Request.Method,
			RelativePath: c.Request.URL.Path,
			UserAgent:    c.Request.UserAgent(),
		}

		switch {
		case c.Writer.Status() >= 500:
			logger.Error(requestLog)
		case c.Writer.Status() >= 400:
			logger.Warn(requestLog)
		default:
			logger.Info(requestLog)
		}
	}
}
