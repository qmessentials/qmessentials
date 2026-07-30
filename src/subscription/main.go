package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const serviceName = "subscription"

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"status":  "UP",
		})
	})
	return router
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	slog.Info("starting service", "service", serviceName, "port", port)
	if err := setupRouter().Run(":" + port); err != nil {
		slog.Error("server failed", "service", serviceName, "error", err)
		os.Exit(1)
	}
}
