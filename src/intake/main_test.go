package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, subject string, data []byte) error {
	args := m.Called(ctx, subject, data)
	return args.Error(0)
}

type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) Subscribe(subject string, handler func(data []byte) error) error {
	args := m.Called(subject, handler)
	return args.Error(0)
}

func TestHealthCheck(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	mockPublisher := new(MockPublisher)
	router := setupRouter(nil, mockPublisher)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), "GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"UP"}`, w.Body.String())
}

func TestPostTestResults(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	mockPublisher := new(MockPublisher)
	router := setupRouter(nil, mockPublisher)

	payload := `{"id": "550e8400-e29b-41d4-a716-446655440000", "sampleId": 123}`
	mockPublisher.On("Publish", mock.Anything, "test-results", []byte(payload)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), "POST", "/test-results", strings.NewReader(payload))
	req.Header.Set("X-Internal-Token", "test-secret")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockPublisher.AssertExpectations(t)
}
