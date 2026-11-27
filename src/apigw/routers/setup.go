package routers

import (
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(credentials auth.Credentials, sessions auth.Sessions, users auth.Users) *gin.Engine {
	r := gin.Default()

	webUrl := os.Getenv("WEB_URL")
	if webUrl != "" {
		slog.Info("CORS enabled", "origin", webUrl)
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{webUrl},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	} else {
		slog.Warn("WEB_URL not set, CORS headers will not be sent")
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})
	intakeUrlStr, ok := os.LookupEnv("INTAKE_URL")
	if !ok {
		intakeUrlStr = "/intake"
	}
	intakeUrl, err := url.Parse(intakeUrlStr)
	if err != nil {
		slog.Error("Invalid INTAKE_URL", "url", intakeUrlStr, "error", err)
		os.Exit(1)
	}
	configUrlStr, ok := os.LookupEnv("CONFIG_URL")
	if !ok {
		configUrlStr = "/config"
	}
	configUrl, err := url.Parse(configUrlStr)
	if err != nil {
		slog.Error("Invalid CONFIG_URL", "url", configUrlStr, "error", err)
		os.Exit(1)
	}
	subscriptionUrlStr, ok := os.LookupEnv("SUBSCRIPTION_URL")
	if !ok {
		subscriptionUrlStr = "/subscription"
	}
	subscriptionUrl, err := url.Parse(subscriptionUrlStr)
	if err != nil {
		slog.Error("Invalid SUBSCRIPTION_URL", "url", subscriptionUrlStr, "error", err)
		os.Exit(1)
	}
	apiSharedSecret, ok := os.LookupEnv("API_SHARED_SECRET")
	if !ok {
		slog.Error("API_SHARED_SECRET must be set")
		os.Exit(1)
	}
	api := r.Group("/api")
	registerAuthRoutes(api, credentials, sessions, users)
	registerAdminRoutes(api.Group("/admin"), sessions, users)
	registerProxyRoutes(api, sessions, apiSharedSecret, intakeUrl, configUrl, subscriptionUrl)
	return r
}
