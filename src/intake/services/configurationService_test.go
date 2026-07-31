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

func TestGetProductMetadata(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", request.Method)
		}
		if request.URL.Path != "/products/KE-TST-HE01/metadata" {
			t.Errorf("unexpected path %s", request.URL.Path)
		}
		if request.Header.Get("X-Internal-Token") != "secret" {
			t.Error("expected internal token header")
		}
		return configurationResponse(http.StatusOK, `{
			"partNumber":"KE-TST-HE01",
			"metadataType":"regex",
			"metadataSelector":"^(?<plant>[A-Z]{3})$",
			"definitions":[{"metadataKey":"plant","valueType":"string"}]
		}`), nil
	})}

	service := NewConfigurationServiceHTTP(client, "http://configuration:8081/", "secret")
	metadata, err := service.GetProductMetadata(context.Background(), "KE-TST-HE01")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if metadata.MetadataType == nil || *metadata.MetadataType != "regex" {
		t.Fatalf("unexpected metadata: %+v", metadata)
	}
	if len(metadata.Definitions) != 1 || metadata.Definitions[0].MetadataKey != "plant" {
		t.Fatalf("unexpected definitions: %+v", metadata.Definitions)
	}
}

func TestGetProductMetadataRejectsUnsuccessfulResponse(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return configurationResponse(http.StatusNotFound, ""), nil
	})}

	service := NewConfigurationServiceHTTP(client, "http://configuration:8081", "secret")
	if _, err := service.GetProductMetadata(context.Background(), "missing"); err == nil {
		t.Fatal("expected an error")
	}
}

func configurationResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
