package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qmessentials/qmessentials/configuration/repositories"
)

func slogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		latency := time.Since(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				slog.Error("request error", "error", e.Error())
			}
		} else {
			slog.Info("request",
				"status", status,
				"method", method,
				"path", path,
				"query", query,
				"ip", c.ClientIP(),
				"latency", latency,
				"user-agent", c.Request.UserAgent(),
			)
		}
	}
}

func setupRouter(productRepo repositories.ProductRepository) *gin.Engine {
	r := gin.New()
	r.Use(slogMiddleware())
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
	r.Use(func(c *gin.Context) {
		token := c.Request.Header.Get("X-Internal-Token")
		if token != apiSharedSecret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	})
	r.GET("/products/:partNumber", func(c *gin.Context) {
		product, err := productRepo.GetByPartNumber(c.Request.Context(), c.Param("partNumber"))
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	// Add more routes here

	return r
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8081"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		mustGetEnv("DB_USER"),
		mustGetEnv("DB_PASSWORD"),
		mustGetEnv("DB_HOST"),
		mustGetEnv("DB_PORT"),
		mustGetEnv("DB_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("failed to open database connection pool", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing database connection pool")
		if err := db.Close(); err != nil {
			slog.Error("failed to close database smoothly", "error", err)
		}
	}()
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("successfully connected to database")

	productRepo := repositories.NewProductRepositoryPG(db)
	r := setupRouter(productRepo)

	slog.Info("starting server", "port", port)
	if err = r.Run(":" + port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		slog.Error("missing environment variable", "key", key)
		os.Exit(1)
	}
	return value
}
