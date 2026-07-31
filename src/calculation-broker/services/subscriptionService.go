package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qmessentials/qmessentials/calculation-broker/models"
)

type SubscriptionService interface {
	GetActive(ctx context.Context) ([]models.Subscription, error)
}

type SubscriptionServiceHTTP struct {
	client          *http.Client
	baseURL         string
	apiSharedSecret string
}

func NewSubscriptionServiceHTTP(client *http.Client, baseURL string, apiSharedSecret string) *SubscriptionServiceHTTP {
	return &SubscriptionServiceHTTP{
		client:          client,
		baseURL:         strings.TrimRight(baseURL, "/"),
		apiSharedSecret: apiSharedSecret,
	}
}

func (s *SubscriptionServiceHTTP) GetActive(ctx context.Context) ([]models.Subscription, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/subscriptions/active", nil)
	if err != nil {
		return nil, fmt.Errorf("create active subscriptions request: %w", err)
	}
	request.Header.Set("X-Internal-Token", s.apiSharedSecret)

	response, err := s.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get active subscriptions: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get active subscriptions: unexpected status %s", response.Status)
	}

	var subscriptions []models.Subscription
	if err = json.NewDecoder(response.Body).Decode(&subscriptions); err != nil {
		return nil, fmt.Errorf("decode active subscriptions: %w", err)
	}
	return subscriptions, nil
}
