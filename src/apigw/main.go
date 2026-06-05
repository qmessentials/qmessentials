package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

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
		panic("Invalid INTAKE_URL: " + intakeUrlStr)
	}
	configUrlStr, ok := os.LookupEnv("CONFIG_URL")
	if !ok {
		configUrlStr = "/config"
	}
	configUrl, err := url.Parse(configUrlStr)
	if err != nil {
		panic("Invalid CONFIG_URL: " + configUrlStr)
	}
	apiSharedSecret, ok := os.LookupEnv("API_SHARED_SECRET")
	if !ok {
		panic("API_SHARED_SECRET must be set")
	}
	r.Any("/api/intake/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, apiSharedSecret, intakeUrl)
	})
	r.Any("/api/config/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, apiSharedSecret, configUrl)
	})
	return r
}

func proxyRequest(c *gin.Context, apiSharedSecret string, apiUrl *url.URL) {
	user, err := getAuthenticatedRequestUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	proxyPath := c.Param("proxyPath")
	proxy := httputil.NewSingleHostReverseProxy(apiUrl)
	proxy.Rewrite = func(req *httputil.ProxyRequest) {
		req.SetXForwarded()
		req.Out.URL.Scheme = apiUrl.Scheme
		req.Out.URL.Host = apiUrl.Host
		req.Out.URL.Path = path.Join(apiUrl.Path, proxyPath)
		req.Out.URL.RawQuery = c.Request.URL.RawQuery
		req.Out.Header.Set("X-Internal-Token", apiSharedSecret)
		req.Out.Header.Set("X-User-Id", user.Name)
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	r := setupRouter()
	if err := r.Run(":" + port); err != nil {
		panic(err)
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
	return "", nil
}

func getUser(userName string) (*UserInfo, error) {
	//TODO: get this from a database
	return &UserInfo{
		ID:    1,
		Name:  userName,
		Roles: []string{},
	}, nil
}
