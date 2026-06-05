package repositories

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)
import "github.com/qmessentials/qmessentials/intake/models"

type SampleRepository interface {
	Get() ([]models.Sample, error)
}

type SampleRepositoryPG struct {
	host     string
	port     string
	database string
	user     string
	password string
}

func NewSampleRepositoryPG(host string, port string, database string, user string, password string) *SampleRepositoryPG {
	return &SampleRepositoryPG{host, port, database, user, password}
}

func (r *SampleRepositoryPG) Get() ([]models.Sample, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", r.user, r.password, r.host, r.port, r.database)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Warn("failed to close database", "error", err)
		}
	}()
	results := make([]models.Sample, 0)
	rows, err := db.Query("select id, serial_number, part_number, status, created_at, updated_at from samples")
	if err != nil {
		slog.Error("failed to query samples", "error", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Warn("failed to close rows", "error", err)
		}
	}()
	for rows.Next() {
		var sample models.Sample
		err = rows.Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt)
		if err != nil {
			slog.Error("failed to scan sample row", "error", err)
			return nil, err
		}
		results = append(results, sample)
	}
	return results, err
}
