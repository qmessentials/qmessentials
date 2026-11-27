package routers

import (
	"errors"
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(admin *gin.RouterGroup, sessions auth.Sessions, users auth.Users) {
	admin.Use(func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		id, err := sessions.Validate(c.Request.Context(), token)
		if err != nil {
			sessionError(c, err)
			c.Abort()
			return
		}
		allowed, err := users.IsAdministrator(c.Request.Context(), id)
		if err != nil {
			slog.Error("Failed to check Administrator role", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Administrator role required"})
			return
		}
		c.Next()
	})
	admin.GET("/users", func(c *gin.Context) {
		list, err := users.List(c.Request.Context())
		if err != nil {
			userError(c, err)
			return
		}
		c.JSON(http.StatusOK, list)
	})
	admin.POST("/users", func(c *gin.Context) {
		var request struct {
			UserID   string   `json:"userId" binding:"required"`
			Password string   `json:"password" binding:"required"`
			Roles    []string `json:"roles"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID, password, and roles are required"})
			return
		}
		if err := users.Create(c.Request.Context(), request.UserID, request.Password, request.Roles); err != nil {
			userError(c, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"userId": strings.ToLower(request.UserID)})
	})
	admin.PUT("/users/:userId/roles", func(c *gin.Context) {
		var request struct {
			Roles []string `json:"roles"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roles"})
			return
		}
		if err := users.SetRoles(c.Request.Context(), c.Param("userId"), request.Roles); err != nil {
			userError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
	admin.PUT("/users/:userId/disabled", func(c *gin.Context) {
		var request struct {
			Disabled *bool `json:"disabled" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil || request.Disabled == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Disabled is required"})
			return
		}
		if err := users.SetDisabled(c.Request.Context(), c.Param("userId"), *request.Disabled); err != nil {
			userError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
}

func userError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidUserID), errors.Is(err, auth.ErrInvalidPassword), errors.Is(err, auth.ErrInvalidRoles):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, auth.ErrUserExists), errors.Is(err, auth.ErrLastAdministrator):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, auth.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		slog.Error("Failed to manage user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
}
