package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupRouter(credentials credentialService, sessions sessionService) *gin.Engine {
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
		var request struct {
			UserID   string `json:"userId" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and password are required"})
			return
		}
		userID, err := credentials.authenticate(c.Request.Context(), request.UserID, request.Password)
		if err != nil {
			if errors.Is(err, errInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID or password"})
				return
			}
			slog.Error("Failed to authenticate user", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		token, expires, err := sessions.create(c.Request.Context(), userID)
		if err != nil {
			slog.Error("Failed to create session", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		setSessionCookie(c, token, expires)
		c.JSON(http.StatusOK, gin.H{"userId": userID, "expiresAt": expires})
	})
	r.POST("/api/session/renew", func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, nextToken, expires, err := sessions.renew(c.Request.Context(), token)
		if err != nil {
			sessionError(c, err)
			return
		}
		setSessionCookie(c, nextToken, expires)
		c.JSON(http.StatusOK, gin.H{"userId": userID, "expiresAt": expires})
	})
	r.GET("/api/session", func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID, err := sessions.validate(c.Request.Context(), token)
		if err != nil {
			sessionError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"userId": userID})
	})
	r.POST("/api/logout", func(c *gin.Context) {
		token, err := sessionToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if err := sessions.revoke(c.Request.Context(), token); err != nil {
			sessionError(c, err)
			return
		}
		clearSessionCookie(c)
		c.Status(http.StatusNoContent)
	})
	r.POST("/api/change-password", func(c *gin.Context) {
		var request struct {
			UserID      string `json:"userId" binding:"required"`
			OldPassword string `json:"oldPassword" binding:"required"`
			NewPassword string `json:"newPassword" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID, current password, and new password are required"})
			return
		}
		err := credentials.changePassword(c.Request.Context(), request.UserID, request.OldPassword, request.NewPassword)
		if errors.Is(err, errInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID or password"})
			return
		}
		if errors.Is(err, errInvalidPassword) || errors.Is(err, errPasswordUnchanged) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			slog.Error("Failed to change password", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Status(http.StatusNoContent)
	})
	r.Any("/api/intake/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, sessions, apiSharedSecret, intakeUrl)
	})
	r.Any("/api/config/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, sessions, apiSharedSecret, configUrl)
	})
	r.Any("/api/subscription/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, sessions, apiSharedSecret, subscriptionUrl)
	})
	return r
}

func proxyRequest(c *gin.Context, sessions sessionService, apiSharedSecret string, apiUrl *url.URL) {
	token, err := sessionToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := sessions.validate(c.Request.Context(), token)
	if err != nil {
		sessionError(c, err)
		return
	}
	proxyPath := c.Param("proxyPath")

	slog.Info("proxying request", "path", proxyPath, "user", userID, "dest", apiUrl.Host)
	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetXForwarded()
			req.Out.URL.Scheme = apiUrl.Scheme
			req.Out.URL.Host = apiUrl.Host
			req.Out.URL.Path = path.Join(apiUrl.Path, proxyPath)
			req.Out.URL.RawQuery = c.Request.URL.RawQuery
			req.Out.Header.Del("Cookie")
			req.Out.Header.Del("Authorization")
			req.Out.Header.Set("X-Internal-Token", apiSharedSecret)
			req.Out.Header.Set("X-User-Id", userID)
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

	db, err := openAuthDB()
	if err != nil {
		slog.Error("auth database unavailable", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	credentials := credentialService{repository: sqlCredentialRepository{db: db}}
	r := setupRouter(credentials, sessionService{repository: sqlSessionRepository{db: db}})
	if err := r.Run(":" + port); err != nil {
		slog.Error("server run failed", "error", err)
		os.Exit(1)
	}
}

func openAuthDB() (*sql.DB, error) {
	for _, name := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		if os.Getenv(name) == "" {
			return nil, fmt.Errorf("%s must be set", name)
		}
	}
	address := &url.URL{
		Scheme: "postgres", Host: os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path: os.Getenv("DB_NAME"),
		User: url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
	}
	query := address.Query()
	query.Set("sslmode", "disable")
	address.RawQuery = query.Encode()
	db, err := sql.Open("pgx", address.String())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func sessionToken(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", errInvalidSession
	}
	return cookie.Value, nil
}

func sessionError(c *gin.Context, err error) {
	if errors.Is(err, errInvalidSession) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	slog.Error("Failed to access session", "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}

func setSessionCookie(c *gin.Context, token string, expires time.Time) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: sessionCookieName, Value: token, Path: "/api", Expires: expires,
		HttpOnly: true, Secure: secureSessionCookie(c), SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name: sessionCookieName, Path: "/api", MaxAge: -1,
		HttpOnly: true, Secure: secureSessionCookie(c), SameSite: http.SameSiteLaxMode,
	})
}

func secureSessionCookie(c *gin.Context) bool {
	return c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
