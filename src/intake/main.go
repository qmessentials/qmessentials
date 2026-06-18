package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/qmessentials/qmessentials/intake/messaging"
	"github.com/qmessentials/qmessentials/intake/repositories"
	"github.com/qmessentials/qmessentials/intake/routers"
	"github.com/qmessentials/qmessentials/intake/services"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8082"
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

	db, closeDB, err := connectDB(ctx, dsn)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer closeDB()
	slog.Info("successfully connected to database")

	sampleRepo := repositories.NewSampleRepositoryPG(db)
	testResultRepo := repositories.NewTestResultRepositoryPG(db)

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, closeNATS, err := connectNATS(natsURL)
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer closeNATS()
	slog.Info("successfully connected to NATS", "url", natsURL)

	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("failed to create JetStream client", "error", err)
		os.Exit(1)
	}

	if _, err = js.Stream(ctx, "TEST_RESULTS"); err != nil {
		if !errors.Is(err, jetstream.ErrStreamNotFound) {
			slog.Error("failed to check JetStream stream", "error", err)
			os.Exit(1)
		}
		if _, err = js.CreateStream(ctx, jetstream.StreamConfig{
			Name:      "TEST_RESULTS",
			Subjects:  []string{"test-results"},
			Retention: jetstream.WorkQueuePolicy,
		}); err != nil {
			slog.Error("failed to create JetStream stream", "error", err)
			os.Exit(1)
		}
	}

	publisher := messaging.NewNatsPublisher(js)
	subscriber := messaging.NewNatsSubscriber(js)

	testResultConsumer := services.NewTestResultConsumer(subscriber, testResultRepo, new(os.Getenv("HASH_KEY")), services.HashPayload)
	if err = testResultConsumer.Subscribe(ctx, "TEST_RESULTS", "test-results"); err != nil {
		slog.Error("failed to subscribe to NATS queue", "error", err)
		os.Exit(1)
	}

	r := routers.Setup(sampleRepo, publisher)

	slog.Info("starting server", "port", port)
	if err = r.Run(":" + port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func connectDB(ctx context.Context, dsn string) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, err
	}
	return pool, pool.Close, nil
}

func connectNATS(url string) (*nats.Conn, func(), error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, nil, err
	}
	return nc, nc.Close, nil
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		slog.Error("missing environment variable", "key", key)
		os.Exit(1)
	}
	return value
}
