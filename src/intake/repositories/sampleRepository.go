// Package repositories provides database access
package repositories

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qmessentials/qmessentials/intake/models"
)

type SampleRepository interface {
	Get(ctx context.Context) ([]models.Sample, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*models.Sample, error)
	Create(ctx context.Context, serialNumber, partNumber string) (*models.Sample, error)
	Transition(ctx context.Context, serialNumber, from, to string) (*models.Sample, error)
}

var ErrSampleNotFound = errors.New("sample not found")
var ErrDuplicateSerial = errors.New("serial number already exists")
var ErrTransitionConflict = errors.New("sample status has changed")

type SampleRepositoryPG struct {
	db *pgxpool.Pool
}

func NewSampleRepositoryPG(db *pgxpool.Pool) *SampleRepositoryPG {
	return &SampleRepositoryPG{db}
}

func (r *SampleRepositoryPG) Get(ctx context.Context) ([]models.Sample, error) {
	results := make([]models.Sample, 0)
	rows, err := r.db.Query(ctx, "select id, serial_number, part_number, status, created_at, updated_at from samples")
	if err != nil {
		slog.Error("failed to query samples", "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sample models.Sample
		if err = rows.Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt); err != nil {
			slog.Error("failed to scan sample row", "error", err)
			return nil, err
		}
		results = append(results, sample)
	}
	return results, rows.Err()
}

func (r *SampleRepositoryPG) GetBySerialNumber(ctx context.Context, serialNumber string) (*models.Sample, error) {
	sample, err := r.getSampleBySerialNumber(ctx, serialNumber)
	if err != nil {
		return nil, err
	}
	testResults, err := r.getTestResultsForSample(ctx, sample.SerialNumber)
	if err != nil {
		return nil, err
	}
	sample.TestResults = &testResults
	return sample, nil
}

func (r *SampleRepositoryPG) getSampleBySerialNumber(ctx context.Context, serialNumber string) (*models.Sample, error) {
	row := r.db.QueryRow(ctx, "select id, serial_number, part_number, status, created_at, updated_at from samples where serial_number = $1", serialNumber)
	var sample models.Sample
	if err := row.Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSampleNotFound
		}
		slog.Error("failed to scan sample row", "error", err)
		return nil, err
	}
	return &sample, nil
}

func (r *SampleRepositoryPG) Create(ctx context.Context, serialNumber, partNumber string) (*models.Sample, error) {
	var sample models.Sample
	err := r.db.QueryRow(ctx, `insert into samples (serial_number, part_number, status)
		values ($1, $2, 'TESTING') returning id, serial_number, part_number, status, created_at, updated_at`, serialNumber, partNumber).
		Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, ErrDuplicateSerial
	}
	return &sample, err
}

func (r *SampleRepositoryPG) Transition(ctx context.Context, serialNumber, from, to string) (*models.Sample, error) {
	var sample models.Sample
	err := r.db.QueryRow(ctx, `update samples set status = $3, updated_at = now()
		where serial_number = $1 and status = $2
		returning id, serial_number, part_number, status, created_at, updated_at`, serialNumber, from, to).
		Scan(&sample.ID, &sample.SerialNumber, &sample.PartNumber, &sample.Status, &sample.CreatedAt, &sample.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTransitionConflict
	}
	return &sample, err
}

func (r *SampleRepositoryPG) getTestResultsForSample(ctx context.Context, serialNumber string) ([]models.TestResult, error) {
	results := make([]models.TestResult, 0)
	rows, err := r.db.Query(ctx, "select id, serial_number, part_number, canonical_test_name, coalesce(modifiers, '{}'::text[]), test_result, unit, decimal_places, min_value, max_value, hash_value, voided_at, voided_by, voided_reason, void_comment, created_at, updated_at from test_results where serial_number = $1", serialNumber)
	if err != nil {
		slog.Error("failed to query test results", "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var result models.TestResult
		if err = rows.Scan(
			&result.ID, &result.SerialNumber, &result.PartNumber, &result.CanonicalTestName,
			&result.Modifiers,
			&result.TestResult, &result.Unit, &result.DecimalPlaces, &result.MinValue, &result.MaxValue,
			&result.HashValue, &result.VoidedAt, &result.VoidedBy, &result.VoidedReason, &result.VoidComment,
			&result.CreatedAt, &result.UpdatedAt,
		); err != nil {
			slog.Error("failed to scan test result row", "error", err)
			return nil, err
		}
		results = append(results, result)
	}
	return results, rows.Err()
}
