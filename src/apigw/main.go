package main

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
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
	r.POST("/api/login", func(c *gin.Context) {
		c.SetCookie("token", "temp-token", 0, "/", "", false, true)
		user, err := getUser("temp-user")
		if err != nil {
			slog.Error("Failed to authenticate user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, user)
	})
	r.Any("/api/intake/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, apiSharedSecret, intakeUrl)
	})
	r.Any("/api/config/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, apiSharedSecret, configUrl)
	})
	r.Any("/api/subscription/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, apiSharedSecret, subscriptionUrl)
	})
	return r
}

func proxyRequest(c *gin.Context, apiSharedSecret string, apiUrl *url.URL) {
	user, err := getAuthenticatedRequestUser(c)
	if err != nil {
		slog.Warn("unauthorized request", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	proxyPath := c.Param("proxyPath")

	slog.Info("proxying request", "path", proxyPath, "user", user.Name, "dest", apiUrl.Host)
	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetXForwarded()
			req.Out.URL.Scheme = apiUrl.Scheme
			req.Out.URL.Host = apiUrl.Host
			req.Out.URL.Path = path.Join(apiUrl.Path, proxyPath)
			req.Out.URL.RawQuery = c.Request.URL.RawQuery
			req.Out.Header.Set("X-Internal-Token", apiSharedSecret)
			req.Out.Header.Set("X-User-Id", user.Name)
		},
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	slog.Info("starting api gateway", "port", port)

	r := setupRouter()
	if err := r.Run(":" + port); err != nil {
		slog.Error("server run failed", "error", err)
		os.Exit(1)
	}
}

type UserInfo struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func getAuthenticatedRequestUser(c *gin.Context) (*UserInfo, error) {
	token, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	userName, err := getUserNameFromToken(token)
	if err != nil {
		return nil, err
	}
	return getUser(userName)
}

func getUserNameFromToken(token string) (string, error) {
	if token == "temp-token" {
		return "temp-user", nil
	}
	return "", errors.New("invalid token")
}

func getUser(userName string) (*UserInfo, error) {
	//TODO: get this from a database
	if userName == "temp-user" {
		return &UserInfo{
			ID:    1,
			Name:  userName,
			Roles: []string{"temp-role"},
		}, nil
	}
	return nil, errors.New("user not found")
}
