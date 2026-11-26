package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userId"

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				slog.Error("request error", "error", err.Error())
			}
			return
		}

		slog.Info("request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"latency", time.Since(start),
			"user-agent", c.Request.UserAgent(),
		)
	}
}

func SharedSecret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Internal-Token") != secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}

func AuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing user identity"})
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}
