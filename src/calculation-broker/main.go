package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qmessentials/qmessentials/calculation-broker/services"
)

const serviceName = "calculation-broker"

func run(ctx context.Context, subscriptionService services.SubscriptionService) error {
	slog.Info("worker started", "service", serviceName)
	subscriptions, err := subscriptionService.GetActive(ctx)
	if err != nil {
		return err
	}
	slog.Info("loaded active subscriptions", "service", serviceName, "count", len(subscriptions))

	<-ctx.Done()
	slog.Info("worker stopped", "service", serviceName)
	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	subscriptionService := services.NewSubscriptionServiceHTTP(
		&http.Client{Timeout: 10 * time.Second},
		mustGetEnv("SUBSCRIPTION_URL"),
		mustGetEnv("API_SHARED_SECRET"),
	)
	if err := run(ctx, subscriptionService); err != nil {
		slog.Error("worker failed", "service", serviceName, "error", err)
		os.Exit(1)
	}
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		slog.Error("missing environment variable", "key", key)
		os.Exit(1)
	}
	return value
}
