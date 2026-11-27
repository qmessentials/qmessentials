package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type sqlCredentialRepository struct{ db *sql.DB }

func (r sqlCredentialRepository) get(ctx context.Context, id string) (credentialRecord, error) {
	var record credentialRecord
	err := r.db.QueryRowContext(ctx, `select user_id, password_hash, disabled from users where user_id = $1`, id).
		Scan(&record.UserID, &record.PasswordHash, &record.Disabled)
	if errors.Is(err, sql.ErrNoRows) {
		return credentialRecord{}, errUserNotFound
	}
	return record, err
}

func (r sqlCredentialRepository) create(ctx context.Context, id, hash string) error {
	_, err := r.db.ExecContext(ctx, `insert into users (user_id, password_hash) values ($1, $2)`, id, hash)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return errUserExists
	}
	return err
}

func (r sqlCredentialRepository) setPassword(ctx context.Context, id, hash string) error {
	result, err := r.db.ExecContext(ctx,
		`update users set password_hash = $2, updated_at = now() where user_id = $1`, id, hash)
	return affectedUser(result, err)
}

func (r sqlCredentialRepository) setDisabled(ctx context.Context, id string, disabled bool) error {
	result, err := r.db.ExecContext(ctx,
		`update users set disabled = $2, updated_at = now() where user_id = $1`, id, disabled)
	return affectedUser(result, err)
}

func affectedUser(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errUserNotFound
	}
	return nil
}
