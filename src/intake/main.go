package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/qmessentials/qmessentials/intake/queue"
	"github.com/qmessentials/qmessentials/intake/repositories"
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

func setupRouter(sampleRepo repositories.SampleRepository, publisher queue.Publisher) *gin.Engine {
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

		c.Status(http.StatusCreated)
	})

	// Add more routes here

	return r
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8082"
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

	sampleRepo := repositories.NewSampleRepositoryPG(db)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Close()
	slog.Info("successfully connected to NATS", "url", natsURL)

	publisher := queue.NewNatsPublisher(nc)
	subscriber := queue.NewNatsSubscriber(nc)

	err = subscriber.Subscribe("test-results", func(data []byte) error {
		slog.Info("received test result from queue", "data", string(data))
		return nil
	})
	if err != nil {
		slog.Error("failed to subscribe to NATS queue", "error", err)
		os.Exit(1)
	}

	r := setupRouter(sampleRepo, publisher)

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
