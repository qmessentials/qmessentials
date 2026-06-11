package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
)
import "github.com/qmessentials/qmessentials/intake/models"

type SampleRepository interface {
	Get(ctx context.Context) ([]models.Sample, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (models.Sample, error)
}

type SampleRepositoryPG struct {
	db *sql.DB
}

func NewSampleRepositoryPG(db *sql.DB) *SampleRepositoryPG {
	return &SampleRepositoryPG{db}
}

func (r *SampleRepositoryPG) Get(ctx context.Context) ([]models.Sample, error) {
	results := make([]models.Sample, 0)
	rows, err := r.db.QueryContext(ctx, "select id, serial_number, part_number, status, created_at, updated_at from samples")
	if err != nil {
		slog.Error("failed to query samples", "error", err)
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
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

func (r *SampleRepositoryPG) GetBySerialNumber(ctx context.Context, serialNumber string) (models.Sample, error) {
	row := r.db.QueryRowContext(ctx, "select id, serial_number, part_number, status, created_at, updated_at from samples where serial_number = $1", serialNumber)
	var sample models.Sample
	err := row.Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt)
	if err != nil {
		slog.Error("failed to scan sample row", "error", err)
		return models.Sample{}, err
	}
	return sample, nil
}
