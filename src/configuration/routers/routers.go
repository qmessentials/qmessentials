// Package routers provides routing services for the HTTP server.
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qmessentials/qmessentials/configuration/middleware"
	"github.com/qmessentials/qmessentials/configuration/repositories"
)

func Setup(productRepo repositories.ProductRepository, apiSharedSecret string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logging())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	internal := r.Group("/")
	internal.Use(middleware.SharedSecret(apiSharedSecret))
	internal.GET("/products/:partNumber/metadata", func(c *gin.Context) {
		metadata, err := productRepo.GetMetadataByPartNumber(c.Request.Context(), c.Param("partNumber"))
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product metadata"})
			return
		}
		c.JSON(http.StatusOK, metadata)
	})

	authorized := internal.Group("/")
	authorized.Use(middleware.AuthenticatedUser())
	authorized.GET("/products/:partNumber", func(c *gin.Context) {
		product, err := productRepo.GetByPartNumber(c.Request.Context(), c.Param("partNumber"))
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	return r
}
