package routers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qmessentials/qmessentials/subscription/models"
	"github.com/qmessentials/qmessentials/subscription/repositories"
)

type subscriptionRepositoryStub struct {
	subscriptions []models.Subscription
	created       *models.Subscription
	err           error
	ownerUserID   string
	activeOnly    bool
	createInput   models.CreateSubscriptionInput
	updateInput   models.UpdateSubscriptionInput
	id            int
	called        bool
}

func (r *subscriptionRepositoryStub) GetForUser(_ context.Context, ownerUserID string, activeOnly bool) ([]models.Subscription, error) {
	r.called = true
	r.ownerUserID = ownerUserID
	r.activeOnly = activeOnly
	return r.subscriptions, r.err
}

func (r *subscriptionRepositoryStub) GetByID(_ context.Context, id int) (*models.Subscription, error) {
	r.called = true
	r.id = id
	return r.created, r.err
}

func (r *subscriptionRepositoryStub) Create(_ context.Context, ownerUserID string, input models.CreateSubscriptionInput) (*models.Subscription, error) {
	r.called = true
	r.ownerUserID = ownerUserID
	r.createInput = input
	return r.created, r.err
}

func (r *subscriptionRepositoryStub) Update(_ context.Context, id int, input models.UpdateSubscriptionInput) (*models.Subscription, error) {
	r.called = true
	r.id = id
	r.updateInput = input
	return r.created, r.err
}

func (r *subscriptionRepositoryStub) Deactivate(_ context.Context, id int) error {
	r.called = true
	r.id = id
	return r.err
}

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	response := performRequest(Setup(&subscriptionRepositoryStub{}, "secret"), "/health", nil)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if !strings.Contains(response.Body.String(), `"status":"UP"`) {
		t.Fatalf("expected UP response, got %s", response.Body.String())
	}
}

func TestGetSubscriptionsDefaultsToActiveOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{subscriptions: []models.Subscription{{
		ID:          7,
		OwnerUserID: "user-123",
		RuleText:    "plant:CHA",
		VersionID:   3,
		IsActive:    true,
		CreatedAt:   time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC),
	}}}
	headers := map[string]string{"X-Internal-Token": "secret", "X-User-Id": "user-123"}

	response := performRequest(Setup(repo, "secret"), "/subscriptions", headers)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if !repo.called || repo.ownerUserID != "user-123" || !repo.activeOnly {
		t.Fatalf("unexpected repository call: called=%v user=%q activeOnly=%v", repo.called, repo.ownerUserID, repo.activeOnly)
	}
	if !strings.Contains(response.Body.String(), `"versionId":3`) {
		t.Fatalf("expected subscription response, got %s", response.Body.String())
	}
}

func TestGetSubscriptionsCanIncludeInactive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{subscriptions: make([]models.Subscription, 0)}
	headers := map[string]string{"X-Internal-Token": "secret", "X-User-Id": "user-123"}

	response := performRequest(Setup(repo, "secret"), "/subscriptions?activeOnly=false", headers)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if repo.activeOnly {
		t.Fatal("expected activeOnly to be false")
	}
	if response.Body.String() != "[]" {
		t.Fatalf("expected empty JSON array, got %s", response.Body.String())
	}
}

func TestGetSubscriptionsRejectsInvalidActiveOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{}
	headers := map[string]string{"X-Internal-Token": "secret", "X-User-Id": "user-123"}

	response := performRequest(Setup(repo, "secret"), "/subscriptions?activeOnly=sometimes", headers)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
	if repo.called {
		t.Fatal("repository should not be called")
	}
}

func TestGetSubscriptionsRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{}

	response := performRequest(Setup(repo, "secret"), "/subscriptions", map[string]string{"X-Internal-Token": "secret"})

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestGetSubscriptionsHandlesRepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{err: errors.New("database unavailable")}
	headers := map[string]string{"X-Internal-Token": "secret", "X-User-Id": "user-123"}

	response := performRequest(Setup(repo, "secret"), "/subscriptions", headers)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestGetSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{created: &models.Subscription{
		ID: 8, OwnerUserID: "user-123", RuleText: "plant:CHA", VersionID: 2, IsActive: true,
	}}

	response := performRequest(Setup(repo, "secret"), "/subscriptions/8", authenticatedJSONHeaders())

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if repo.id != 8 {
		t.Fatalf("expected subscription ID 8, got %d", repo.id)
	}
	if !strings.Contains(response.Body.String(), `"ruleText":"plant:CHA"`) {
		t.Fatalf("expected subscription response, got %s", response.Body.String())
	}
}

func TestGetSubscriptionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{err: repositories.ErrSubscriptionNotFound}

	response := performRequest(Setup(repo, "secret"), "/subscriptions/99", authenticatedJSONHeaders())

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestCreateSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{created: &models.Subscription{
		ID:          8,
		OwnerUserID: "user-123",
		RuleText:    "plant:CHA",
		VersionID:   1,
		IsActive:    true,
	}}
	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Internal-Token": "secret",
		"X-User-Id":        "user-123",
	}

	response := performRequestWithBody(
		Setup(repo, "secret"),
		http.MethodPost,
		"/subscriptions",
		`{"ruleText":"plant:CHA"}`,
		headers,
	)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}
	if repo.ownerUserID != "user-123" || repo.createInput.RuleText != "plant:CHA" {
		t.Fatalf("unexpected repository call: user=%q input=%+v", repo.ownerUserID, repo.createInput)
	}
	if !strings.Contains(response.Body.String(), `"versionId":1`) {
		t.Fatalf("expected created subscription response, got %s", response.Body.String())
	}
}

func TestCreateSubscriptionRequiresRuleText(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{}
	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Internal-Token": "secret",
		"X-User-Id":        "user-123",
	}

	response := performRequestWithBody(Setup(repo, "secret"), http.MethodPost, "/subscriptions", `{"ruleText":"  "}`, headers)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
	if repo.called {
		t.Fatal("repository should not be called")
	}
}

func TestCreateSubscriptionHandlesRepositoryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{err: errors.New("database unavailable")}
	headers := map[string]string{
		"Content-Type":     "application/json",
		"X-Internal-Token": "secret",
		"X-User-Id":        "user-123",
	}

	response := performRequestWithBody(
		Setup(repo, "secret"),
		http.MethodPost,
		"/subscriptions",
		`{"ruleText":"plant:CHA"}`,
		headers,
	)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}

func TestUpdateSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{created: &models.Subscription{
		ID: 8, OwnerUserID: "user-123", RuleText: "plant:MEM", VersionID: 2, IsActive: true,
	}}
	headers := authenticatedJSONHeaders()

	response := performRequestWithBody(
		Setup(repo, "secret"), http.MethodPut, "/subscriptions/8", `{"ruleText":"plant:MEM"}`, headers,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if repo.id != 8 || repo.updateInput.RuleText != "plant:MEM" {
		t.Fatalf("unexpected repository call: id=%d input=%+v", repo.id, repo.updateInput)
	}
	if !strings.Contains(response.Body.String(), `"versionId":2`) {
		t.Fatalf("expected updated subscription response, got %s", response.Body.String())
	}
}

func TestUpdateSubscriptionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{err: repositories.ErrSubscriptionNotFound}

	response := performRequestWithBody(
		Setup(repo, "secret"), http.MethodPut, "/subscriptions/99", `{"ruleText":"plant:MEM"}`, authenticatedJSONHeaders(),
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestDeleteSubscription(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{}

	response := performRequestWithBody(
		Setup(repo, "secret"), http.MethodDelete, "/subscriptions/8", "", authenticatedJSONHeaders(),
	)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
	if repo.id != 8 {
		t.Fatalf("unexpected repository call: id=%d", repo.id)
	}
}

func TestDeleteSubscriptionNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &subscriptionRepositoryStub{err: repositories.ErrSubscriptionNotFound}

	response := performRequestWithBody(
		Setup(repo, "secret"), http.MethodDelete, "/subscriptions/99", "", authenticatedJSONHeaders(),
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func authenticatedJSONHeaders() map[string]string {
	return map[string]string{
		"Content-Type":     "application/json",
		"X-Internal-Token": "secret",
		"X-User-Id":        "user-123",
	}
}

func performRequest(handler http.Handler, target string, headers map[string]string) *httptest.ResponseRecorder {
	return performRequestWithBody(handler, http.MethodGet, target, "", headers)
}

func performRequestWithBody(handler http.Handler, method string, target string, body string, headers map[string]string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
