package repositories

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qmessentials/qmessentials/intake/models"
)

type TestResultRepository interface {
	Add(ctx context.Context, item *models.TestResult) (uuid.UUID, error)
	GetOne(ctx context.Context, id uuid.UUID) (*models.TestResult, error)
	Update(ctx context.Context, props TestResultUpdateProps) error
}

type TestResultRepositoryPG struct {
	db *pgxpool.Pool
}

func NewTestResultRepositoryPG(db *pgxpool.Pool) *TestResultRepositoryPG {
	return &TestResultRepositoryPG{db}
}

func (r *TestResultRepositoryPG) GetOne(ctx context.Context, id uuid.UUID) (*models.TestResult, error) {
	var result models.TestResult
	err := r.db.QueryRow(ctx,
		`select id, serial_number, part_number, canonical_test_name, modifiers, test_result, unit, decimal_places, min_value, max_value, hash_value, voided_at, voided_by, voided_reason, void_comment, created_at, updated_at from test_results where id = $1`,
		id,
	).Scan(
		&result.ID,
		&result.SerialNumber,
		&result.PartNumber,
		&result.CanonicalTestName,
		&result.Modifiers,
		&result.TestResult,
		&result.Unit,
		&result.DecimalPlaces,
		&result.MinValue,
		&result.MaxValue,
		&result.HashValue,
		&result.VoidedAt,
		&result.VoidedBy,
		&result.VoidedReason,
		&result.VoidComment,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *TestResultRepositoryPG) Add(ctx context.Context, item *models.TestResult) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		`insert into test_results
		(serial_number, part_number, canonical_test_name, modifiers, test_result, unit, decimal_places, min_value, max_value, hash_value)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id`,
		item.SerialNumber,
		item.PartNumber,
		item.CanonicalTestName,
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

type TestResultUpdateProps struct {
	Id           uuid.UUID
	VoidedAt     *time.Time
	VoidedBy     *string
	VoidedReason *string
	VoidComment  *string
}

func (r *TestResultRepositoryPG) Update(ctx context.Context, props TestResultUpdateProps) error {
	_, err := r.db.Exec(ctx,
		`update test_results
		set voided_at = $1, voided_by = $2, voided_reason = $3, void_comment = $4
		where id = $5`,
		props.VoidedAt,
		props.VoidedBy,
		props.VoidedReason,
		props.VoidComment,
		props.Id,
	)
	if err != nil {
		slog.Error("failed to update test result", "error", err)
		return err
	}
	return nil
}
