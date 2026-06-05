package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/qmessentials/qmessentials/intake/repositories"
)

func setupRouter(sampleRepo repositories.SampleRepository) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	apiSharedSecret, ok := os.LookupEnv("API_SHARED_SECRET")
	if !ok {
		panic("API_SHARED_SECRET must be set")
	}
	r.Use(func(c *gin.Context) {
		token := c.Request.Header.Get("X-Internal-Token")
		if token != apiSharedSecret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	})

	r.GET("/samples", func(c *gin.Context) {
		samples, err := sampleRepo.Get()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
			return
		}
		c.JSON(http.StatusOK, samples)
	})

	// Add more routes here

	return r
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8082"
	}
	sampleRepo := repositories.NewSampleRepositoryPG(mustGetEnv("DB_HOST"), mustGetEnv("DB_PORT"), mustGetEnv("DB_NAME"), mustGetEnv("DB_USER"), mustGetEnv("DB_PASSWORD"))
	r := setupRouter(sampleRepo)
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(key + " must be set")
	}
	return value
}
