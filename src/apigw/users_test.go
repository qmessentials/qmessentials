package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"github.com/qmessentials/qmessentials/apigw/routers"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

type memoryUserRepository struct{ users map[string]auth.User }

func (r *memoryUserRepository) List(_ context.Context) ([]auth.User, error) {
	list := []auth.User{}
	for _, u := range r.users {
		list = append(list, u)
	}
	return list, nil
}
func (r *memoryUserRepository) Get(_ context.Context, id string) (auth.User, error) {
	u, ok := r.users[id]
	if !ok {
		return auth.User{}, auth.ErrUserNotFound
	}
	return u, nil
}
func (r *memoryUserRepository) Create(_ context.Context, id, _ string, roles []string) error {
	if _, exists := r.users[id]; exists {
		return auth.ErrUserExists
	}
	r.users[id] = auth.User{UserID: id, Roles: roles}
	return nil
}
func (r *memoryUserRepository) SetDisabled(_ context.Context, id string, disabled bool) error {
	u, ok := r.users[id]
	if !ok {
		return auth.ErrUserNotFound
	}
	if disabled && slices.Contains(u.Roles, auth.AdministratorRole) && r.enabledAdministratorCount(id) == 0 {
		return auth.ErrLastAdministrator
	}
	u.Disabled = disabled
	r.users[id] = u
	return nil
}
func (r *memoryUserRepository) SetRoles(_ context.Context, id string, roles []string) error {
	u, ok := r.users[id]
	if !ok {
		return auth.ErrUserNotFound
	}
	if !u.Disabled && slices.Contains(u.Roles, auth.AdministratorRole) && !slices.Contains(roles, auth.AdministratorRole) && r.enabledAdministratorCount(id) == 0 {
		return auth.ErrLastAdministrator
	}
	u.Roles = roles
	r.users[id] = u
	return nil
}
func (r *memoryUserRepository) enabledAdministratorCount(exclude string) int {
	count := 0
	for id, u := range r.users {
		if id != exclude && !u.Disabled && slices.Contains(u.Roles, auth.AdministratorRole) {
			count++
		}
	}
	return count
}
func (r *memoryUserRepository) Bootstrap(_ context.Context, id, _ string) error {
	for _, u := range r.users {
		if slices.Contains(u.Roles, auth.AdministratorRole) {
			return nil
		}
	}
	if _, exists := r.users[id]; exists {
		return auth.ErrUserExists
	}
	r.users[id] = auth.User{UserID: id, Roles: []string{auth.AdministratorRole}}
	return nil
}

func testUsers(ids ...string) auth.Users {
	repository := &memoryUserRepository{users: map[string]auth.User{}}
	for _, id := range ids {
		repository.users[id] = auth.User{UserID: id, Roles: []string{}}
	}
	return auth.NewUsers(repository)
}

func TestUserRoleValidation(t *testing.T) {
	repository := &memoryUserRepository{users: map[string]auth.User{}}
	service := auth.NewUsers(repository)
	for _, roles := range [][]string{nil, {"Owner"}, {auth.AdministratorRole, auth.AdministratorRole}} {
		if err := service.Create(context.Background(), "person", "valid-password-123", roles); !errors.Is(err, auth.ErrInvalidRoles) {
			t.Fatalf("expected invalid roles for %v, got %v", roles, err)
		}
	}
	if err := service.Create(context.Background(), "unassigned", "valid-password-123", []string{}); err != nil {
		t.Fatal(err)
	}
	if len(repository.users["unassigned"].Roles) != 0 {
		t.Fatal("explicit empty role list was rejected")
	}
	roles := []string{auth.TechnicianRole, auth.AdministratorRole}
	if err := service.Create(context.Background(), "person", "valid-password-123", roles); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(repository.users["person"].Roles, roles) {
		t.Fatalf("create changed role order: %v", repository.users["person"].Roles)
	}
}

func TestUserManagementHTTP(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "test-secret")
	usersRepo := &memoryUserRepository{users: map[string]auth.User{}}
	users := auth.NewUsers(usersRepo)
	if err := users.Bootstrap(context.Background(), "admin", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	if err := users.Bootstrap(context.Background(), "other", "valid-password-123"); err != nil {
		t.Fatal(err)
	}
	if len(usersRepo.users) != 1 {
		t.Fatal("bootstrap created a second account")
	}
	credentialsRepo := &memoryCredentialRepository{records: map[string]auth.CredentialRecord{}}
	credentials := auth.NewCredentials(credentialsRepo)
	for _, id := range []string{"admin", "technician"} {
		if err := credentials.Create(context.Background(), id, "valid-password-123"); err != nil {
			t.Fatal(err)
		}
	}
	usersRepo.users["technician"] = auth.User{UserID: "technician", Roles: []string{auth.TechnicianRole}}
	router := routers.Setup(credentials, auth.NewSessions(newMemorySessionRepository()), users)
	call := func(method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		if cookie != nil {
			req.AddCookie(cookie)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		return response
	}
	login := func(id string) *http.Cookie {
		response := call("POST", "/api/login", `{"userId":"`+id+`","password":"valid-password-123"}`, nil)
		if response.Code != 200 {
			t.Fatalf("login %s: %d", id, response.Code)
		}
		return response.Result().Cookies()[0]
	}
	admin, technician := login("admin"), login("technician")
	if response := call("GET", "/api/admin/users", "", nil); response.Code != 401 {
		t.Fatalf("anonymous: %d", response.Code)
	}
	if response := call("GET", "/api/admin/users", "", technician); response.Code != 403 {
		t.Fatalf("nonadmin: %d", response.Code)
	}
	if response := call("POST", "/api/admin/users", `{"userId":"newuser","password":"valid-password-123","roles":["Technician"]}`, technician); response.Code != 403 {
		t.Fatalf("nonadmin create: %d", response.Code)
	}
	if response := call("POST", "/api/admin/users", `{"userId":"newuser","password":"valid-password-123","roles":["Technician"]}`, admin); response.Code != 201 {
		t.Fatalf("create: %d %s", response.Code, response.Body.String())
	}
	if response := call("PUT", "/api/admin/users/technician/roles", `{"roles":["Administrator","Technician"]}`, admin); response.Code != 204 {
		t.Fatalf("assign: %d", response.Code)
	}
	if response := call("GET", "/api/admin/users", "", technician); response.Code != 200 {
		t.Fatalf("role not effective: %d", response.Code)
	}
	if response := call("PUT", "/api/admin/users/technician/roles", `{"roles":["Technician"]}`, admin); response.Code != 204 {
		t.Fatalf("remove: %d", response.Code)
	}
	if response := call("GET", "/api/admin/users", "", technician); response.Code != 403 {
		t.Fatalf("removed role still effective: %d", response.Code)
	}
	if response := call("PUT", "/api/admin/users/admin/disabled", `{"disabled":true}`, admin); response.Code != 409 {
		t.Fatalf("last admin disabled: %d", response.Code)
	}
	if response := call("PUT", "/api/admin/users/admin/roles", `{"roles":[]}`, admin); response.Code != 409 {
		t.Fatalf("last admin role removed: %d", response.Code)
	}
	if response := call("POST", "/api/admin/users", `{"userId":"baduser","password":"valid-password-123","roles":["Owner"]}`, admin); response.Code != 400 {
		t.Fatalf("unknown role: %d", response.Code)
	}
	if err := users.SetDisabled(context.Background(), "missing", true); !errors.Is(err, auth.ErrUserNotFound) {
		t.Fatal(err)
	}
}
