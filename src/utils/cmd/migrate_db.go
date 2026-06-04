package cmd

import (
	"cmp"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateDbCmd = &cobra.Command{
	Use:   "migrate-db",
	Short: "Run database migrations",
	Long:  `Executes database migrations for Postgres using the specified connection parameters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := viper.GetString("host")
		port := viper.GetInt("port")
		database := viper.GetString("database")
		user := viper.GetString("user")
		password := viper.GetString("password")
		path := viper.GetString("path")

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, database)
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			return fmt.Errorf("failed to open database connection: %w", err)
		}
		defer func() {
			err := db.Close()
			if err != nil {
				cmd.Println("failed to close database connection:", err)
			}
		}()

		domain := filepath.Base(path)
		cmd.Printf("Checking schema migrations for domain '%s'\n", domain)

		ctx := cmd.Context()

		var currentVersion int64
		var exists bool
		err = db.QueryRowContext(ctx, "select to_regclass('schema_migrations') is not null").Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check if schema_migrations table exists: %w", err)
		}

		if !exists {
			currentVersion = -1
		} else {
			query := "select coalesce(max(version_number), -1) from schema_migrations where domain = $1"
			err = db.QueryRowContext(ctx, query, domain).Scan(&currentVersion)
			if err != nil {
				return fmt.Errorf("failed to retrieve current migration version: %w", err)
			}
		}
		cmd.Printf("Current migration version for domain '%s': %d\n", domain, currentVersion)

		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("failed to read migrations directory: %w", err)
		}

		type migration struct {
			version int64
			name    string
		}
		var migrations []migration

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasSuffix(name, ".sql") {
				continue
			}

			parts := strings.SplitN(name, "-", 2)
			if len(parts) < 2 {
				continue
			}

			v, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				continue
			}

			if v > currentVersion {
				migrations = append(migrations, migration{version: v, name: name})
			}
		}

		slices.SortFunc(migrations, func(a, b migration) int {
			return cmp.Compare(a.version, b.version)
		})

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				cmd.Println("failed to rollback transaction:", err)
			}
		}()

		for _, m := range migrations {
			cmd.Printf("Applying migration %s for domain '%s'\n", m.name, domain)
			content, err := os.ReadFile(filepath.Join(path, m.name))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", m.name, err)
			}

			_, err = tx.ExecContext(ctx, string(content))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", m.name, err)
			}

			_, err = tx.ExecContext(ctx, "insert into schema_migrations (domain, version_number) values ($1, $2)", domain, m.version)
			if err != nil {
				return fmt.Errorf("failed to update schema_migrations for %s: %w", m.name, err)
			}
		}
		cmd.Printf("Migration completed for domain '%s'\n", domain)
		return tx.Commit()
	},
}

func init() {
	rootCmd.AddCommand(migrateDbCmd)

	migrateDbCmd.Flags().String("host", "localhost", "Database host")
	migrateDbCmd.Flags().Int("port", 5432, "Database port")
	migrateDbCmd.Flags().String("database", "postgres", "Database name")
	migrateDbCmd.Flags().String("user", "postgres", "Database user")
	migrateDbCmd.Flags().String("password", "", "Database password")
	migrateDbCmd.Flags().StringP("path", "p", "", "Path to migration files")

	mustBindPFlag("host", migrateDbCmd)
	mustBindPFlag("port", migrateDbCmd)
	mustBindPFlag("database", migrateDbCmd)
	mustBindPFlag("user", migrateDbCmd)
	mustBindPFlag("password", migrateDbCmd)
	mustBindPFlag("path", migrateDbCmd)

	err := migrateDbCmd.MarkFlagRequired("path")
	if err != nil {
		panic("failed to mark path flag as required")
	}
}
