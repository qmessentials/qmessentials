package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"github.com/qmessentials/qmessentials/apigw/routers"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type memoryCredentialRepository struct {
	records map[string]auth.CredentialRecord
}

func (r *memoryCredentialRepository) Get(_ context.Context, id string) (auth.CredentialRecord, error) {
	record, ok := r.records[id]
	if !ok {
		return auth.CredentialRecord{}, auth.ErrUserNotFound
	}
	return record, nil
}
func (r *memoryCredentialRepository) Create(_ context.Context, id, hash string) error {
	if _, ok := r.records[id]; ok {
		return auth.ErrUserExists
	}
	r.records[id] = auth.CredentialRecord{UserID: id, PasswordHash: hash}
	return nil
}
func (r *memoryCredentialRepository) SetPassword(_ context.Context, id, hash string) error {
	record, ok := r.records[id]
	if !ok {
		return auth.ErrUserNotFound
	}
	record.PasswordHash = hash
	r.records[id] = record
	return nil
}
func (r *memoryCredentialRepository) SetDisabled(_ context.Context, id string, disabled bool) error {
	record, ok := r.records[id]
	if !ok {
		return auth.ErrUserNotFound
	}
	record.Disabled = disabled
	r.records[id] = record
	return nil
}

func TestCredentialLifecycle(t *testing.T) {
	ctx := context.Background()
	repository := &memoryCredentialRepository{records: map[string]auth.CredentialRecord{}}
	service := auth.NewCredentials(repository)
	if err := service.Create(ctx, "Alice", "first-password-123"); err != nil {
		t.Fatal(err)
	}
	if repository.records["alice"].PasswordHash == "first-password-123" {
		t.Fatal("plaintext password stored")
	}
	storedCost, err := bcrypt.Cost([]byte(repository.records["alice"].PasswordHash))
	if err != nil || storedCost != 12 {
		t.Fatalf("stored bcrypt cost: %d, %v", storedCost, err)
	}
	if err := service.Create(ctx, "ALICE", "another-password-123"); !errors.Is(err, auth.ErrUserExists) {
		t.Fatalf("duplicate: %v", err)
	}
	if id, err := service.Authenticate(ctx, "ALICE", "first-password-123"); err != nil || id != "alice" {
		t.Fatalf("valid login: %q, %v", id, err)
	}
	for _, tc := range []struct{ id, password string }{
		{"alice", "wrong-password"}, {"unknown", "first-password-123"}, {"bad id", "first-password-123"},
	} {
		if _, err := service.Authenticate(ctx, tc.id, tc.password); !errors.Is(err, auth.ErrInvalidCredentials) {
			t.Fatalf("failed login %+v: %v", tc, err)
		}
	}
	if err := service.ChangePassword(ctx, "alice", "wrong-password", "second-password-123"); !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("wrong current password: %v", err)
	}
	if err := service.ChangePassword(ctx, "alice", "first-password-123", "second-password-123"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Authenticate(ctx, "alice", "first-password-123"); !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("old password: %v", err)
	}
	if err := service.ResetPassword(ctx, "alice", "reset-password-123"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Authenticate(ctx, "alice", "reset-password-123"); err != nil {
		t.Fatal(err)
	}
	if err := service.SetDisabled(ctx, "alice", true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Authenticate(ctx, "alice", "reset-password-123"); !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("disabled: %v", err)
	}
	if err := service.SetDisabled(ctx, "alice", false); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Authenticate(ctx, "alice", "reset-password-123"); err != nil {
		t.Fatal(err)
	}
}

func TestLoginHTTP(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	repository := &memoryCredentialRepository{records: map[string]auth.CredentialRecord{}}
	service := auth.NewCredentials(repository)
	if err := service.Create(context.Background(), "alice", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	router := routers.Setup(service, auth.NewSessions(newMemorySessionRepository()), testUsers("alice"))
	for _, tc := range []struct {
		body   string
		status int
		text   string
	}{
		{`{"userId":"alice","password":"valid-password-123"}`, 200, `"userId":"alice"`},
		{`{"userId":"alice","password":"wrong-password"}`, 401, "Invalid user ID or password"},
		{`{"userId":"missing","password":"valid-password-123"}`, 401, "Invalid user ID or password"},
		{`{}`, 400, "required"},
	} {
		request := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(tc.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != tc.status || !bytes.Contains(response.Body.Bytes(), []byte(tc.text)) {
			t.Fatalf("%s: %d %s", tc.body, response.Code, response.Body.String())
		}
		if tc.status == http.StatusOK && !bytes.Contains([]byte(response.Header().Get("Set-Cookie")), []byte("HttpOnly")) {
			t.Fatal("login did not issue an HttpOnly session cookie")
		}
		if bytes.Contains(response.Body.Bytes(), []byte("valid-password-123")) {
			t.Fatal("password returned in response")
		}
	}
}

func TestChangePasswordHTTP(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	repository := &memoryCredentialRepository{records: map[string]auth.CredentialRecord{}}
	service := auth.NewCredentials(repository)
	if err := service.Create(context.Background(), "alice", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	router := routers.Setup(service, auth.NewSessions(newMemorySessionRepository()), testUsers("alice"))
	for _, tc := range []struct {
		body   string
		status int
	}{
		{`{"userId":"alice","oldPassword":"wrong-password","newPassword":"another-password-123"}`, 401},
		{`{"userId":"alice","oldPassword":"valid-password-123","newPassword":"short"}`, 400},
		{`{"userId":"alice","oldPassword":"valid-password-123","newPassword":"another-password-123"}`, 204},
	} {
		request := httptest.NewRequest(http.MethodPost, "/api/change-password", bytes.NewBufferString(tc.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != tc.status {
			t.Fatalf("%s: %d %s", tc.body, response.Code, response.Body.String())
		}
	}
	if _, err := service.Authenticate(context.Background(), "alice", "another-password-123"); err != nil {
		t.Fatal(err)
	}
}

func TestTemporaryTokenCannotAccessServices(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	router := routers.Setup(auth.NewCredentials(&memoryCredentialRepository{records: map[string]auth.CredentialRecord{}}), auth.NewSessions(newMemorySessionRepository()), testUsers())
	request := httptest.NewRequest(http.MethodGet, "/api/config/products/example", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: "temp-token"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("temporary token accepted: %d %s", response.Code, response.Body.String())
	}
}
