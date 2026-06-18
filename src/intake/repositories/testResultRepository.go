package repositories

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/qmessentials/qmessentials/intake/models"
)

type TestResultRepository interface {
	Add(ctx context.Context, item *models.TestResult) (uuid.UUID, error)
}

type TestResultRepositoryPG struct {
	db *sql.DB
}

func NewTestResultRepositoryPG(db *sql.DB) *TestResultRepositoryPG {
	return &TestResultRepositoryPG{db}
}

func (r *TestResultRepositoryPG) Add(ctx context.Context, item *models.TestResult) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx,
		`insert into test_results
		(sample_id, part_number, product_test_sequence, modifiers, test_result, unit, decimal_places, min_value, max_value, hash_value)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id`,
		item.SampleID,
		item.PartNumber,
		item.ProductTestSequence,
		item.Modifiers,
		item.TestResult,
		item.Unit,
		item.DecimalPlaces,
		item.MinValue,
		item.MaxValue,
		item.HashValue,
	).Scan(&id)
	if err != nil {
		slog.Error("failed to insert test result", "error", err)
		return uuid.Nil, err
	}
	return id, nil
}
