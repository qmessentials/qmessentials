package services

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestGetActiveSubscriptions(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", request.Method)
		}
		if request.URL.String() != "http://subscription:8083/subscriptions/active" {
			t.Errorf("unexpected URL %s", request.URL)
		}
		if request.Header.Get("X-Internal-Token") != "secret" {
			t.Error("expected internal token header")
		}
		return response(http.StatusOK, `[{"id":7,"ruleText":"plant:CHA","versionId":3,"isActive":true}]`), nil
	})}

	service := NewSubscriptionServiceHTTP(client, "http://subscription:8083/", "secret")
	subscriptions, err := service.GetActive(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(subscriptions) != 1 || subscriptions[0].ID != 7 || subscriptions[0].RuleText != "plant:CHA" {
		t.Fatalf("unexpected subscriptions: %+v", subscriptions)
	}
}

func TestGetActiveSubscriptionsRejectsUnsuccessfulResponse(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusServiceUnavailable, ""), nil
	})}

	service := NewSubscriptionServiceHTTP(client, "http://subscription:8083", "secret")
	if _, err := service.GetActive(context.Background()); err == nil {
		t.Fatal("expected an error")
	}
}

func TestGetActiveSubscriptionsRejectsInvalidJSON(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "not-json"), nil
	})}

	service := NewSubscriptionServiceHTTP(client, "http://subscription:8083", "secret")
	if _, err := service.GetActive(context.Background()); err == nil {
		t.Fatal("expected an error")
	}
}

func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
