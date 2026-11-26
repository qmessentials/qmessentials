// Package routers provides routing services for the HTTP server
package routers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qmessentials/qmessentials/intake/messaging"
	"github.com/qmessentials/qmessentials/intake/middleware"
	"github.com/qmessentials/qmessentials/intake/repositories"
	"github.com/qmessentials/qmessentials/intake/services"
)

func Setup(sampleRepo repositories.SampleRepository, testResultRepo repositories.TestResultRepository, publisher messaging.Publisher, configuration services.ConfigurationService) *gin.Engine {
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
			if errors.Is(err, repositories.ErrSampleNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
				return
			}
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sample"})
			return
		}
		c.JSON(http.StatusOK, sample)
	})
	r.POST("/samples", func(c *gin.Context) {
		if !hasRole(c, "Technician") && !hasRole(c, "Quality Manager") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Technician permission required"})
			return
		}
		var input struct {
			SerialNumber string `json:"serialNumber"`
			PartNumber   string `json:"partNumber"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.SerialNumber) == "" || strings.TrimSpace(input.PartNumber) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Serial number and product are required"})
			return
		}
		input.SerialNumber = strings.TrimSpace(input.SerialNumber)
		input.PartNumber = strings.TrimSpace(input.PartNumber)
		if _, err := configuration.GetProductMetadata(c.Request.Context(), input.PartNumber); err != nil {
			if errors.Is(err, services.ErrProductNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Product does not exist"})
				return
			}
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Could not validate product"})
			return
		}
		sample, err := sampleRepo.Create(c.Request.Context(), input.SerialNumber, input.PartNumber)
		if errors.Is(err, repositories.ErrDuplicateSerial) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Serial number already exists"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not create sample"})
			return
		}
		c.JSON(http.StatusCreated, sample)
	})
	r.PATCH("/samples/:serialNumber/status", func(c *gin.Context) {
		var input struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
			return
		}
		current, err := sampleRepo.GetBySerialNumber(c.Request.Context(), c.Param("serialNumber"))
		if errors.Is(err, repositories.ErrSampleNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch sample"})
			return
		}
		allowed := (current.Status == "TESTING" && input.Status == "COMPLETED" && (hasRole(c, "Technician") || hasRole(c, "Quality Manager"))) ||
			(current.Status == "COMPLETED" && (input.Status == "TESTING" || input.Status == "ARCHIVED") && hasRole(c, "Quality Manager"))
		if !allowed {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Status transition is not allowed for this user"})
			return
		}
		sample, err := sampleRepo.Transition(c.Request.Context(), current.SerialNumber, current.Status, input.Status)
		if errors.Is(err, repositories.ErrTransitionConflict) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Sample status changed; refresh and retry"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not update sample"})
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
		var input struct {
			SerialNumber string `json:"serialNumber"`
		}
		if err := json.Unmarshal(body, &input); err != nil || input.SerialNumber == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Serial number is required"})
			return
		}
		sample, err := sampleRepo.GetBySerialNumber(c.Request.Context(), input.SerialNumber)
		if errors.Is(err, repositories.ErrSampleNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch sample"})
			return
		}
		if sample.Status != "TESTING" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Sample is not available for testing"})
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

func hasRole(c *gin.Context, role string) bool {
	if c.GetHeader("X-User-Id") == "" {
		return false
	}
	for _, value := range strings.Split(c.GetHeader("X-User-Roles"), ",") {
		if strings.TrimSpace(value) == role {
			return true
		}
	}
	return false
}
