package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projdocs/cli/pkg/federation/response"
	"go.uber.org/zap"
)

// Recovery returns a Gin middleware that catches panics, logs them with
// a full stack trace via zap, and writes a 500 response to the client.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if val := recover(); val != nil {
				log.Error("panic recovered",
					zap.Any("error", val),
					zap.Stack("stack"),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

func NoRoute(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Warn("no route",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusNotFound, response.Failure(
			"NOT_FOUND",
			"the requested resource does not exist",
		))
	}
}

func NoMethod(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Warn("method not allowed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusMethodNotAllowed, response.Failure(
			"METHOD_NOT_ALLOWED",
			"method not allowed on this resource",
		))
	}
}
