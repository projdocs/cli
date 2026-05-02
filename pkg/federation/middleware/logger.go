package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger returns a Gin middleware that emits one structured log
// line per request, including method, path, status, latency, and client IP.
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if q := c.Request.URL.RawQuery; q != "" {
			path += "?" + q
		}

		c.Next() // execute the handler chain

		log.Debug("request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
