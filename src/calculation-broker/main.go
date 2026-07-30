package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const serviceName = "calculation-broker"

type subscriptionEvent struct {
	SubscriptionID string `json:"subscriptionId" binding:"required"`
	Revision       int64  `json:"revision" binding:"required"`
	Action         string `json:"action" binding:"required"`
}

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"status":  "UP",
		})
	})
	router.POST("/subscriptions/events", func(c *gin.Context) {
		var event subscriptionEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		slog.Info(
			"received subscription event",
			"subscriptionId", event.SubscriptionID,
			"revision", event.Revision,
			"action", event.Action,
		)
		c.Status(http.StatusAccepted)
	})
	return router
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	slog.Info("starting service", "service", serviceName, "port", port)
	if err := setupRouter().Run(":" + port); err != nil {
		slog.Error("server failed", "service", serviceName, "error", err)
		os.Exit(1)
	}
}
