package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/qmessentials/qmessentials/intake/models"
)

type ConfigurationService interface {
	GetProductMetadata(ctx context.Context, partNumber string) (*models.ProductMetadata, error)
}

type ConfigurationServiceHTTP struct {
	client          *http.Client
	baseURL         string
	apiSharedSecret string
}

func NewConfigurationServiceHTTP(client *http.Client, baseURL string, apiSharedSecret string) *ConfigurationServiceHTTP {
	return &ConfigurationServiceHTTP{
		client:          client,
		baseURL:         strings.TrimRight(baseURL, "/"),
		apiSharedSecret: apiSharedSecret,
	}
}

func (s *ConfigurationServiceHTTP) GetProductMetadata(ctx context.Context, partNumber string) (*models.ProductMetadata, error) {
	requestURL := fmt.Sprintf("%s/products/%s/metadata", s.baseURL, url.PathEscape(partNumber))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create product metadata request: %w", err)
	}
	request.Header.Set("X-Internal-Token", s.apiSharedSecret)

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get product metadata: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get product metadata: unexpected status %s", response.Status)
	}

	var metadata models.ProductMetadata
	if err = json.NewDecoder(response.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("decode product metadata: %w", err)
	}
	return &metadata, nil
}
