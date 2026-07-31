package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/qmessentials/qmessentials/configuration/repositories"
	"github.com/qmessentials/qmessentials/configuration/routers"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8081"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		mustGetEnv("DB_USER"),
		mustGetEnv("DB_PASSWORD"),
		mustGetEnv("DB_HOST"),
		mustGetEnv("DB_PORT"),
		mustGetEnv("DB_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("failed to open database connection pool", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing database connection pool")
		if err := db.Close(); err != nil {
			slog.Error("failed to close database smoothly", "error", err)
		}
	}()
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("successfully connected to database")

	productRepo := repositories.NewProductRepositoryPG(db)
	r := routers.Setup(productRepo, mustGetEnv("API_SHARED_SECRET"))

	slog.Info("starting server", "port", port)
	if err = r.Run(":" + port); err != nil {
		slog.Error("server failed", "error", err)
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
