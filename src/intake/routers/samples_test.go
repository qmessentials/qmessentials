package routers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qmessentials/qmessentials/intake/models"
	"github.com/qmessentials/qmessentials/intake/repositories"
	"github.com/qmessentials/qmessentials/intake/services"
)

type sampleStore struct{ sample *models.Sample }

func (s *sampleStore) Get(context.Context) ([]models.Sample, error) { return nil, nil }
func (s *sampleStore) GetBySerialNumber(_ context.Context, serial string) (*models.Sample, error) {
	if s.sample == nil || s.sample.SerialNumber != serial {
		return nil, repositories.ErrSampleNotFound
	}
	return s.sample, nil
}
func (s *sampleStore) Create(_ context.Context, serial, part string) (*models.Sample, error) {
	if s.sample != nil && s.sample.SerialNumber == serial {
		return nil, repositories.ErrDuplicateSerial
	}
	s.sample = &models.Sample{SerialNumber: serial, PartNumber: part, Status: "TESTING"}
	return s.sample, nil
}
func (s *sampleStore) Transition(_ context.Context, serial, from, to string) (*models.Sample, error) {
	if s.sample == nil || s.sample.SerialNumber != serial || s.sample.Status != from {
		return nil, repositories.ErrTransitionConflict
	}
	s.sample.Status = to
	return s.sample, nil
}

type productLookup struct{}

func (productLookup) GetProductMetadata(_ context.Context, part string) (*models.ProductMetadata, error) {
	if part != "P1" {
		return nil, services.ErrProductNotFound
	}
	return &models.ProductMetadata{}, nil
}

type testPublisher struct{ calls int }

func (p *testPublisher) Publish(context.Context, string, []byte) error { p.calls++; return nil }

func TestSampleLifecycleAndResultGate(t *testing.T) {
	t.Setenv("API_SHARED_SECRET", "secret")
	store := &sampleStore{}
	publisher := &testPublisher{}
	router := Setup(store, nil, publisher, productLookup{})
	request := func(method, path, body, role string) int {
		t.Helper()
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("X-Internal-Token", "secret")
		req.Header.Set("X-User-Id", "user")
		req.Header.Set("X-User-Roles", role)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		return response.Code
	}
	if got := request("POST", "/samples", `{"serialNumber":"S1","partNumber":"P1"}`, "Technician"); got != http.StatusCreated {
		t.Fatalf("create: %d", got)
	}
	if store.sample.Status != "TESTING" {
		t.Fatalf("new status: %s", store.sample.Status)
	}
	if got := request("POST", "/samples", `{"serialNumber":"S1","partNumber":"P1"}`, "Technician"); got != http.StatusConflict {
		t.Fatalf("duplicate: %d", got)
	}
	if got := request("POST", "/samples", `{"serialNumber":"S2","partNumber":"missing"}`, "Technician"); got != http.StatusBadRequest {
		t.Fatalf("invalid product: %d", got)
	}
	if got := request("POST", "/test-results", `{"serialNumber":"S1"}`, "Technician"); got != http.StatusAccepted {
		t.Fatalf("open result: %d", got)
	}
	if got := request("PATCH", "/samples/S1/status", `{"status":"COMPLETED"}`, "Technician"); got != http.StatusOK {
		t.Fatalf("complete: %d", got)
	}
	if got := request("POST", "/test-results", `{"serialNumber":"S1"}`, "Technician"); got != http.StatusConflict {
		t.Fatalf("closed result: %d", got)
	}
	if publisher.calls != 1 {
		t.Fatalf("published results: %d", publisher.calls)
	}
	if got := request("PATCH", "/samples/S1/status", `{"status":"TESTING"}`, "Technician"); got != http.StatusConflict {
		t.Fatalf("technician reopen: %d", got)
	}
	if got := request("PATCH", "/samples/S1/status", `{"status":"TESTING"}`, "Quality Manager"); got != http.StatusOK {
		t.Fatalf("manager reopen: %d", got)
	}
	if got := request("PATCH", "/samples/S1/status", `{"status":"COMPLETED"}`, "Technician"); got != http.StatusOK {
		t.Fatalf("complete again: %d", got)
	}
	if got := request("PATCH", "/samples/S1/status", `{"status":"ARCHIVED"}`, "Quality Manager"); got != http.StatusOK {
		t.Fatalf("archive: %d", got)
	}
	if got := request("PATCH", "/samples/S1/status", `{"status":"TESTING"}`, "Quality Manager"); got != http.StatusConflict {
		t.Fatalf("terminal reopen: %d", got)
	}
}
