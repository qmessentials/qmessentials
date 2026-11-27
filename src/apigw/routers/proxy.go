package routers

import (
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
)

func registerProxyRoutes(r *gin.RouterGroup, sessions auth.Sessions, apiSharedSecret string, intakeUrl, configUrl, subscriptionUrl *url.URL) {
	r.Any("/intake/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, sessions, apiSharedSecret, intakeUrl)
	})
	r.Any("/config/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, sessions, apiSharedSecret, configUrl)
	})
	r.Any("/subscription/*proxyPath", func(c *gin.Context) {
		proxyRequest(c, sessions, apiSharedSecret, subscriptionUrl)
	})
}

func proxyRequest(c *gin.Context, sessions auth.Sessions, apiSharedSecret string, apiUrl *url.URL) {
	token, err := sessionToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, err := sessions.Validate(c.Request.Context(), token)
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
