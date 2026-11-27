package routers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/qmessentials/qmessentials/configuration/models"
)

type productRepositoryStub struct {
	metadata   *models.ProductMetadata
	products   []models.Product
	activeOnly *bool
}

func (r productRepositoryStub) ListProducts(_ context.Context, activeOnly bool) ([]models.Product, error) {
	if r.activeOnly != nil {
		*r.activeOnly = activeOnly
	}
	return r.products, nil
}

func TestListProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		name, path string
		want       bool
		status     int
	}{
		{"default", "/products", true, http.StatusOK},
		{"active only", "/products?activeOnly=true", true, http.StatusOK},
		{"all products", "/products?activeOnly=false", false, http.StatusOK},
		{"invalid flag", "/products?activeOnly=maybe", false, http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			var actual bool
			repo := productRepositoryStub{
				products:   []models.Product{{ID: 1, PartNumber: "P-100", ProductName: "Example", IsActive: true}},
				activeOnly: &actual,
			}
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request.Header.Set("X-Internal-Token", "secret")
			request.Header.Set("X-User-Id", "technician")
			response := httptest.NewRecorder()
			Setup(repo, "secret").ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status: got %d, want %d", response.Code, test.status)
			}
			if test.status == http.StatusOK && (actual != test.want || !strings.Contains(response.Body.String(), `"partNumber":"P-100"`)) {
				t.Fatalf("unexpected product list: activeOnly=%v body=%s", actual, response.Body.String())
			}
		})
	}
}

func (r productRepositoryStub) GetByPartNumber(context.Context, string) (*models.Product, error) {
	return nil, nil
}

func (r productRepositoryStub) GetMetadataByPartNumber(context.Context, string) (*models.ProductMetadata, error) {
	return r.metadata, nil
}

func TestGetInternalProductMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metadataType := "regex"
	metadataSelector := "^(?<plant>[A-Z]{3})$"
	repo := productRepositoryStub{metadata: &models.ProductMetadata{
		PartNumber:       "KE-TST-HE01",
		MetadataType:     &metadataType,
		MetadataSelector: &metadataSelector,
		Definitions: []models.ProductMetadataDefinition{{
			MetadataKey: "plant",
			ValueType:   "string",
		}},
	}}
	request := httptest.NewRequest(http.MethodGet, "/products/KE-TST-HE01/metadata", nil)
	request.Header.Set("X-Internal-Token", "secret")
	response := httptest.NewRecorder()

	Setup(repo, "secret").ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"metadataType":"regex"`) {
		t.Fatalf("expected product metadata response, got %s", response.Body.String())
	}
}

func TestInternalProductMetadataRequiresToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest(http.MethodGet, "/products/KE-TST-HE01/metadata", nil)
	response := httptest.NewRecorder()

	Setup(productRepositoryStub{}, "secret").ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestProductRouteRequiresUserIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest(http.MethodGet, "/products/KE-TST-HE01", nil)
	request.Header.Set("X-Internal-Token", "secret")
	response := httptest.NewRecorder()

	Setup(productRepositoryStub{}, "secret").ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}
