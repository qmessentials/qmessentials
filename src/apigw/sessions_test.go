package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"github.com/qmessentials/qmessentials/apigw/routers"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type memorySession struct {
	userID  string
	expires time.Time
	revoked bool
}

type memorySessionRepository struct {
	mu       sync.Mutex
	sessions map[string]memorySession
	disabled map[string]bool
}

func newMemorySessionRepository() *memorySessionRepository {
	return &memorySessionRepository{sessions: make(map[string]memorySession), disabled: make(map[string]bool)}
}

func (r *memorySessionRepository) Create(_ context.Context, hash []byte, userID string, expires time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[string(hash)] = memorySession{userID: userID, expires: expires}
	return nil
}

func (r *memorySessionRepository) Get(_ context.Context, hash []byte, now time.Time) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[string(hash)]
	if !ok || session.revoked || !session.expires.After(now) || r.disabled[session.userID] {
		return "", auth.ErrInvalidSession
	}
	return session.userID, nil
}

func (r *memorySessionRepository) Renew(_ context.Context, oldHash, newHash []byte, now, expires time.Time) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[string(oldHash)]
	if !ok || session.revoked || !session.expires.After(now) || r.disabled[session.userID] {
		return "", auth.ErrInvalidSession
	}
	delete(r.sessions, string(oldHash))
	session.expires = expires
	r.sessions[string(newHash)] = session
	return session.userID, nil
}

func (r *memorySessionRepository) Revoke(_ context.Context, hash []byte, now time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[string(hash)]
	if !ok || session.revoked || !session.expires.After(now) {
		return auth.ErrInvalidSession
	}
	session.revoked = true
	r.sessions[string(hash)] = session
	return nil
}

func TestSessionLifecycle(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC)
	repository := newMemorySessionRepository()
	sessions := auth.NewSessions(repository).WithClock(func() time.Time { return now })
	token, expires, err := sessions.Create(ctx, "alice")
	if err != nil || !expires.Equal(now.Add((12 * time.Hour))) {
		t.Fatalf("create: %v, %v", expires, err)
	}
	if _, ok := repository.sessions[token]; ok {
		t.Fatal("raw token stored")
	}
	if user, err := sessions.Validate(ctx, token); err != nil || user != "alice" {
		t.Fatalf("validate: %q, %v", user, err)
	}
	repository.disabled["alice"] = true
	if _, err := sessions.Validate(ctx, token); !errors.Is(err, auth.ErrInvalidSession) {
		t.Fatalf("disabled user session accepted: %v", err)
	}
	repository.disabled["alice"] = false
	for _, invalid := range []string{"", "garbage", token + "x"} {
		if _, err := sessions.Validate(ctx, invalid); !errors.Is(err, auth.ErrInvalidSession) {
			t.Fatalf("invalid token %q: %v", invalid, err)
		}
	}
	now = now.Add(time.Hour)
	_, renewedToken, renewedExpiry, err := sessions.Renew(ctx, token)
	if err != nil || !renewedExpiry.Equal(now.Add((12 * time.Hour))) || renewedToken == token {
		t.Fatalf("renew: %v, %v", renewedExpiry, err)
	}
	if _, err := sessions.Validate(ctx, token); !errors.Is(err, auth.ErrInvalidSession) {
		t.Fatalf("old token still accepted: %v", err)
	}
	if err := sessions.Revoke(ctx, renewedToken); err != nil {
		t.Fatal(err)
	}
	if _, err := sessions.Validate(ctx, renewedToken); !errors.Is(err, auth.ErrInvalidSession) {
		t.Fatalf("revoked token accepted: %v", err)
	}
	expiringToken, _, err := sessions.Create(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add((12 * time.Hour))
	if _, err := sessions.Validate(ctx, expiringToken); !errors.Is(err, auth.ErrInvalidSession) {
		t.Fatalf("expired token accepted: %v", err)
	}
}

func TestSessionHTTPAndProxy(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	credentials := &memoryCredentialRepository{records: map[string]auth.CredentialRecord{}}
	if err := (auth.NewCredentials(credentials)).Create(context.Background(), "alice", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	repository := newMemorySessionRepository()
	sessions := auth.NewSessions(repository)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-User-Id") != "alice" || r.Header.Get("X-Internal-Token") != "test-secret" || r.Header.Get("Cookie") != "" || r.Header.Get("Authorization") != "" {
			t.Errorf("unsafe proxy headers: %v", r.Header)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer backend.Close()
	t.Setenv("CONFIG_URL", backend.URL)
	router := routers.Setup(auth.NewCredentials(credentials), sessions, testUsers("alice"))

	login := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"userId":"alice","password":"valid-password-123"}`))
	login.Header.Set("Content-Type", "application/json")
	loginResponse := httptest.NewRecorder()
	router.ServeHTTP(loginResponse, login)
	if loginResponse.Code != http.StatusOK || len(loginResponse.Result().Cookies()) != 1 {
		t.Fatalf("login: %d %s", loginResponse.Code, loginResponse.Body.String())
	}
	cookie := loginResponse.Result().Cookies()[0]
	request := func(method, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, nil)
		if cookie != nil {
			req.AddCookie(cookie)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		return response
	}
	if response := request(http.MethodGet, "/api/session", cookie); response.Code != http.StatusOK || !bytes.Contains(response.Body.Bytes(), []byte(`"userId":"alice"`)) {
		t.Fatalf("session status: %d %s", response.Code, response.Body.String())
	}
	gateway := httptest.NewServer(router)
	defer gateway.Close()
	proxyReq, err := http.NewRequest(http.MethodGet, gateway.URL+"/api/config/products/example", nil)
	if err != nil {
		t.Fatal(err)
	}
	proxyReq.AddCookie(cookie)
	proxyReq.Header.Set("Authorization", "Bearer untrusted")
	proxyReq.Header.Set("X-User-Id", "mallory")
	proxyResponse, err := gateway.Client().Do(proxyReq)
	if err != nil {
		t.Fatal(err)
	}
	defer proxyResponse.Body.Close()
	if proxyResponse.StatusCode != http.StatusNoContent {
		t.Fatalf("proxy: %d", proxyResponse.StatusCode)
	}
	renewResponse := request(http.MethodPost, "/api/session/renew", cookie)
	if renewResponse.Code != http.StatusOK || len(renewResponse.Result().Cookies()) != 1 {
		t.Fatalf("renew: %d %s", renewResponse.Code, renewResponse.Body.String())
	}
	newCookie := renewResponse.Result().Cookies()[0]
	if response := request(http.MethodGet, "/api/session", cookie); response.Code != http.StatusUnauthorized {
		t.Fatalf("old token after renewal: %d", response.Code)
	}
	if response := request(http.MethodPost, "/api/logout", cookie); response.Code != http.StatusUnauthorized {
		t.Fatalf("old token after renewal: %d", response.Code)
	}
	if response := request(http.MethodPost, "/api/logout", newCookie); response.Code != http.StatusNoContent {
		t.Fatalf("logout: %d %s", response.Code, response.Body.String())
	}
	if response := request(http.MethodPost, "/api/session/renew", newCookie); response.Code != http.StatusUnauthorized {
		t.Fatalf("revoked token renewed: %d", response.Code)
	}
	if response := request(http.MethodGet, "/api/config/products/example", newCookie); response.Code != http.StatusUnauthorized {
		t.Fatalf("revoked token accessed proxy: %d", response.Code)
	}
}
