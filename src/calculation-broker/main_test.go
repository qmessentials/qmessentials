package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/qmessentials/qmessentials/calculation-broker/models"
)

type subscriptionServiceStub struct {
	subscriptions []models.Subscription
	err           error
}

func (s subscriptionServiceStub) GetActive(context.Context) ([]models.Subscription, error) {
	return s.subscriptions, s.err
}

func TestRunStopsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})

	go func() {
		_ = run(ctx, subscriptionServiceStub{})
		close(stopped)
	}()
	cancel()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("worker did not stop after context cancellation")
	}
}

func TestRunReturnsSubscriptionLoadError(t *testing.T) {
	expected := errors.New("subscription service unavailable")
	err := run(context.Background(), subscriptionServiceStub{err: expected})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}
