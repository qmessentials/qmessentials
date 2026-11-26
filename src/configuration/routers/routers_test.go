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
	metadata *models.ProductMetadata
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
