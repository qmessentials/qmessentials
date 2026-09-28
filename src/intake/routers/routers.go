// Package routers provides routing services for the HTTP server
package routers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qmessentials/qmessentials/intake/messaging"
	"github.com/qmessentials/qmessentials/intake/middleware"
	"github.com/qmessentials/qmessentials/intake/repositories"
)

func Setup(sampleRepo repositories.SampleRepository, testResultRepo repositories.TestResultRepository, publisher messaging.Publisher) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logging())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	apiSharedSecret, ok := os.LookupEnv("API_SHARED_SECRET")
	if !ok {
		slog.Error("API_SHARED_SECRET must be set")
		os.Exit(1)
	}
	r.Use(middleware.SharedSecret(apiSharedSecret))

	r.GET("/samples", func(c *gin.Context) {
		samples, err := sampleRepo.Get(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
			return
		}
		slog.Info("fetched samples", "count", len(samples))
		c.JSON(http.StatusOK, samples)
	})
	r.GET("/samples/:serialNumber", func(c *gin.Context) {
		sample, err := sampleRepo.GetBySerialNumber(c.Request.Context(), c.Param("serialNumber"))
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sample"})
			return
		}
		c.JSON(http.StatusOK, sample)
	})
	r.POST("/test-results", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		if err = publisher.Publish(c.Request.Context(), "test-results", body); err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish test result"})
			return
		}

		c.Status(http.StatusAccepted)
	})
	r.PATCH("/test-results/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid test result id"})
			return
		}
		body, err := c.GetRawData()
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		var updateProps repositories.TestResultUpdateProps
		if err = json.Unmarshal(body, &updateProps); err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal request body"})
			return
		}
		updateProps.Id = id
		err = testResultRepo.Update(c.Request.Context(), updateProps)
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test result"})
			return
		}
		c.Status(http.StatusNoContent)
	})
	return r
}
