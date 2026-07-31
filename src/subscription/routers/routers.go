// Package routers provides routing services for the HTTP server.
package routers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qmessentials/qmessentials/subscription/middleware"
	"github.com/qmessentials/qmessentials/subscription/models"
	"github.com/qmessentials/qmessentials/subscription/repositories"
)

func Setup(subscriptionRepo repositories.SubscriptionRepository, apiSharedSecret string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logging())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "subscription", "status": "UP"})
	})

	authorized := r.Group("/")
	authorized.Use(middleware.SharedSecret(apiSharedSecret))
	authorized.Use(middleware.AuthenticatedUser())
	authorized.GET("/subscriptions", func(c *gin.Context) {
		activeOnly := true
		if value, present := c.GetQuery("activeOnly"); present {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "activeOnly must be true or false"})
				return
			}
			activeOnly = parsed
		}

		subscriptions, err := subscriptionRepo.GetForUser(
			c.Request.Context(),
			c.GetString(middleware.UserIDKey),
			activeOnly,
		)
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
			return
		}
		c.JSON(http.StatusOK, subscriptions)
	})
	authorized.GET("/subscriptions/:id", func(c *gin.Context) {
		id, err := subscriptionID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
			return
		}

		subscription, err := subscriptionRepo.GetByID(c.Request.Context(), id)
		if errors.Is(err, repositories.ErrSubscriptionNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscription"})
			return
		}
		c.JSON(http.StatusOK, subscription)
	})
	authorized.POST("/subscriptions", func(c *gin.Context) {
		var input models.CreateSubscriptionInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription"})
			return
		}
		if strings.TrimSpace(input.RuleText) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ruleText is required"})
			return
		}

		subscription, err := subscriptionRepo.Create(
			c.Request.Context(),
			c.GetString(middleware.UserIDKey),
			input,
		)
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
			return
		}
		c.JSON(http.StatusCreated, subscription)
	})
	authorized.PUT("/subscriptions/:id", func(c *gin.Context) {
		id, err := subscriptionID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
			return
		}

		var input models.UpdateSubscriptionInput
		if err = c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.RuleText) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ruleText is required"})
			return
		}

		subscription, err := subscriptionRepo.Update(
			c.Request.Context(),
			id,
			input,
		)
		if errors.Is(err, repositories.ErrSubscriptionNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}
		c.JSON(http.StatusOK, subscription)
	})
	authorized.DELETE("/subscriptions/:id", func(c *gin.Context) {
		id, err := subscriptionID(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
			return
		}

		err = subscriptionRepo.Deactivate(c.Request.Context(), id)
		if errors.Is(err, repositories.ErrSubscriptionNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	return r
}

func subscriptionID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid subscription ID")
	}
	return id, nil
}
