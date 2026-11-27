package routers

import (
	"errors"
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.RouterGroup, credentials auth.Credentials, sessions auth.Sessions, users auth.Users) {
	r.POST("/login", func(c *gin.Context) {
		var request struct {
			UserID   string `json:"userId" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and password are required"})
			return
		}
		userID, err := credentials.Authenticate(c.Request.Context(), request.UserID, request.Password)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID or password"})
				return
			}
			slog.Error("Failed to authenticate user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		token, expires, err := sessions.Create(c.Request.Context(), userID)
		if err != nil {
			slog.Error("Failed to create session", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		setSessionCookie(c, token, expires)
		c.JSON(http.StatusOK, gin.H{"userId": userID, "expiresAt": expires})
	})
	r.POST("/session/renew", func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, nextToken, expires, err := sessions.Renew(c.Request.Context(), token)
		if err != nil {
			sessionError(c, err)
			return
		}
		setSessionCookie(c, nextToken, expires)
		c.JSON(http.StatusOK, gin.H{"userId": userID, "expiresAt": expires})
	})
	r.GET("/session", func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, err := sessions.Validate(c.Request.Context(), token)
		if err != nil {
			sessionError(c, err)
			return
		}
		u, err := users.Get(c.Request.Context(), userID)
		if err != nil {
			userError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"userId": userID, "roles": u.Roles})
	})
	r.POST("/logout", func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if err := sessions.Revoke(c.Request.Context(), token); err != nil {
			sessionError(c, err)
			return
		}
		clearSessionCookie(c)
		c.Status(http.StatusNoContent)
	})
	r.POST("/change-password", func(c *gin.Context) {
		var request struct {
			UserID      string `json:"userId" binding:"required"`
			OldPassword string `json:"oldPassword" binding:"required"`
			NewPassword string `json:"newPassword" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID, current password, and new password are required"})
			return
		}
		err := credentials.ChangePassword(c.Request.Context(), request.UserID, request.OldPassword, request.NewPassword)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID or password"})
			return
		}
		if errors.Is(err, auth.ErrInvalidPassword) || errors.Is(err, auth.ErrPasswordUnchanged) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			slog.Error("Failed to change password", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Status(http.StatusNoContent)
	})
}

func sessionToken(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie(auth.SessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", auth.ErrInvalidSession
	}
	return cookie.Value, nil
}

func sessionError(c *gin.Context, err error) {
	if errors.Is(err, auth.ErrInvalidSession) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	slog.Error("Failed to access session", "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

func setSessionCookie(c *gin.Context, token string, expires time.Time) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: auth.SessionCookieName, Value: token, Path: "/api", Expires: expires,
		HttpOnly: true, Secure: secureSessionCookie(c), SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: auth.SessionCookieName, Path: "/api", MaxAge: -1,
		HttpOnly: true, Secure: secureSessionCookie(c), SameSite: http.SameSiteLaxMode,
	})
}

func secureSessionCookie(c *gin.Context) bool {
	return c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
