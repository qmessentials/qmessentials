package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "calculation-engine"

func run(ctx context.Context) {
	slog.Info("worker started", "service", serviceName)
	<-ctx.Done()
	slog.Info("worker stopped", "service", serviceName)
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	run(ctx)
}
