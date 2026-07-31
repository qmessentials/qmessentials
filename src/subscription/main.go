package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qmessentials/qmessentials/subscription/repositories"
	"github.com/qmessentials/qmessentials/subscription/routers"
)

const serviceName = "subscription"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		mustGetEnv("DB_USER"),
		mustGetEnv("DB_PASSWORD"),
		mustGetEnv("DB_HOST"),
		mustGetEnv("DB_PORT"),
		mustGetEnv("DB_NAME"),
	)
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	if err = db.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}

	subscriptionRepo := repositories.NewSubscriptionRepositoryPG(db)
	router := routers.Setup(subscriptionRepo, mustGetEnv("API_SHARED_SECRET"))

	slog.Info("starting service", "service", serviceName, "port", port)
	if err = router.Run(":" + port); err != nil {
		slog.Error("server failed", "service", serviceName, "error", err)
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
