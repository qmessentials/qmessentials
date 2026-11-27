package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type memoryCredentialRepository struct{ records map[string]credentialRecord }

func (r *memoryCredentialRepository) get(_ context.Context, id string) (credentialRecord, error) {
	record, ok := r.records[id]
	if !ok {
		return credentialRecord{}, errUserNotFound
	}
	return record, nil
}
func (r *memoryCredentialRepository) create(_ context.Context, id, hash string) error {
	if _, ok := r.records[id]; ok {
		return errUserExists
	}
	r.records[id] = credentialRecord{UserID: id, PasswordHash: hash}
	return nil
}
func (r *memoryCredentialRepository) setPassword(_ context.Context, id, hash string) error {
	record, ok := r.records[id]
	if !ok {
		return errUserNotFound
	}
	record.PasswordHash = hash
	r.records[id] = record
	return nil
}
func (r *memoryCredentialRepository) setDisabled(_ context.Context, id string, disabled bool) error {
	record, ok := r.records[id]
	if !ok {
		return errUserNotFound
	}
	record.Disabled = disabled
	r.records[id] = record
	return nil
}

func TestCredentialLifecycle(t *testing.T) {
	ctx := context.Background()
	repository := &memoryCredentialRepository{records: map[string]credentialRecord{}}
	service := credentialService{repository: repository}
	if err := service.create(ctx, "Alice", "first-password-123"); err != nil {
		t.Fatal(err)
	}
	if repository.records["alice"].PasswordHash == "first-password-123" {
		t.Fatal("plaintext password stored")
	}
	storedCost, err := bcrypt.Cost([]byte(repository.records["alice"].PasswordHash))
	if err != nil || storedCost != bcryptCost {
		t.Fatalf("stored bcrypt cost: %d, %v", storedCost, err)
	}
	dummyCost, err := bcrypt.Cost(dummyHash)
	if err != nil || dummyCost != bcryptCost {
		t.Fatalf("dummy bcrypt cost: %d, %v", dummyCost, err)
	}
	if err := service.create(ctx, "ALICE", "another-password-123"); !errors.Is(err, errUserExists) {
		t.Fatalf("duplicate: %v", err)
	}
	if id, err := service.authenticate(ctx, "ALICE", "first-password-123"); err != nil || id != "alice" {
		t.Fatalf("valid login: %q, %v", id, err)
	}
	for _, tc := range []struct{ id, password string }{
		{"alice", "wrong-password"}, {"unknown", "first-password-123"}, {"bad id", "first-password-123"},
	} {
		if _, err := service.authenticate(ctx, tc.id, tc.password); !errors.Is(err, errInvalidCredentials) {
			t.Fatalf("failed login %+v: %v", tc, err)
		}
	}
	if err := service.changePassword(ctx, "alice", "wrong-password", "second-password-123"); !errors.Is(err, errInvalidCredentials) {
		t.Fatalf("wrong current password: %v", err)
	}
	if err := service.changePassword(ctx, "alice", "first-password-123", "second-password-123"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.authenticate(ctx, "alice", "first-password-123"); !errors.Is(err, errInvalidCredentials) {
		t.Fatalf("old password: %v", err)
	}
	if err := service.resetPassword(ctx, "alice", "reset-password-123"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.authenticate(ctx, "alice", "reset-password-123"); err != nil {
		t.Fatal(err)
	}
	if err := service.setDisabled(ctx, "alice", true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.authenticate(ctx, "alice", "reset-password-123"); !errors.Is(err, errInvalidCredentials) {
		t.Fatalf("disabled: %v", err)
	}
	if err := service.setDisabled(ctx, "alice", false); err != nil {
		t.Fatal(err)
	}
	if _, err := service.authenticate(ctx, "alice", "reset-password-123"); err != nil {
		t.Fatal(err)
	}
}

func TestLoginHTTP(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	repository := &memoryCredentialRepository{records: map[string]credentialRecord{}}
	service := credentialService{repository: repository}
	if err := service.create(context.Background(), "alice", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	router := setupRouter(service)
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
		if response.Header().Get("Set-Cookie") != "" {
			t.Fatal("login issued a session cookie before session support exists")
		}
		if bytes.Contains(response.Body.Bytes(), []byte("valid-password-123")) {
			t.Fatal("password returned in response")
		}
	}
}

func TestChangePasswordHTTP(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	repository := &memoryCredentialRepository{records: map[string]credentialRecord{}}
	service := credentialService{repository: repository}
	if err := service.create(context.Background(), "alice", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	router := setupRouter(service)
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
	if _, err := service.authenticate(context.Background(), "alice", "another-password-123"); err != nil {
		t.Fatal(err)
	}
}

func TestTemporaryTokenCannotAccessServices(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	router := setupRouter(credentialService{repository: &memoryCredentialRepository{records: map[string]credentialRecord{}}})
	request := httptest.NewRequest(http.MethodGet, "/api/config/products/example", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: "temp-token"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("temporary token accepted: %d %s", response.Code, response.Body.String())
	}
}
