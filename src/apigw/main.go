package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/qmessentials/qmessentials/apigw/internal/auth"
	"github.com/qmessentials/qmessentials/apigw/routers"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	slog.Info("starting api gateway", "port", port)

	db, err := openAuthDB()
	if err != nil {
		slog.Error("auth database unavailable", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	credentials := auth.NewCredentials(auth.NewSQLCredentialRepository(db))
	users := auth.NewUsers(auth.NewSQLUserRepository(db))
	bootstrapID, bootstrapPassword := os.Getenv("BOOTSTRAP_ADMIN_USER_ID"), os.Getenv("BOOTSTRAP_ADMIN_PASSWORD")
	if (bootstrapID == "") != (bootstrapPassword == "") {
		slog.Error("Both bootstrap Administrator variables must be set together")
		os.Exit(1)
	}
	if bootstrapID != "" {
		if err := users.Bootstrap(context.Background(), bootstrapID, bootstrapPassword); err != nil {
			slog.Error("Failed to bootstrap Administrator", "error", err)
			os.Exit(1)
		}
	}
	r := routers.Setup(credentials, auth.NewSessions(auth.NewSQLSessionRepository(db)), users)
	if err := r.Run(":" + port); err != nil {
		slog.Error("server run failed", "error", err)
		os.Exit(1)
	}
}

func openAuthDB() (*sql.DB, error) {
	for _, name := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		if os.Getenv(name) == "" {
			return nil, fmt.Errorf("%s must be set", name)
		}
	}
	address := &url.URL{
		Scheme: "postgres", Host: os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path: os.Getenv("DB_NAME"),
		User: url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
	}
	query := address.Query()
	query.Set("sslmode", "disable")
	address.RawQuery = query.Encode()
	db, err := sql.Open("pgx", address.String())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
