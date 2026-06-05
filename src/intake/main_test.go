package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	router := setupRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"UP"}`, w.Body.String())
}
