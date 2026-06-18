package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		latency := time.Since(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				slog.Error("request error", "error", e.Error())
			}
		} else {
			slog.Info("request",
				"status", status,
				"method", method,
				"path", path,
				"query", query,
				"ip", c.ClientIP(),
				"latency", latency,
				"user-agent", c.Request.UserAgent(),
			)
		}
	}
}

func SharedSecret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("X-Internal-Token")
		if token != secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}
