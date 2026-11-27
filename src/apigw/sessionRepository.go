package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type sqlSessionRepository struct{ db *sql.DB }

func (r sqlSessionRepository) create(ctx context.Context, hash []byte, userID string, expires time.Time) error {
	_, err := r.db.ExecContext(ctx, `insert into sessions (token_hash, user_id, expires_at) values ($1, $2, $3)`, hash, userID, expires)
	return err
}

func (r sqlSessionRepository) get(ctx context.Context, hash []byte, now time.Time) (string, error) {
	var userID string
	err := r.db.QueryRowContext(ctx, `select s.user_id from sessions s join users u on u.user_id = s.user_id
		where s.token_hash = $1 and s.expires_at > $2 and s.revoked_at is null and not u.disabled`, hash, now).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errInvalidSession
	}
	return userID, err
}

func (r sqlSessionRepository) renew(ctx context.Context, oldHash, newHash []byte, now, expires time.Time) (string, error) {
	var userID string
	err := r.db.QueryRowContext(ctx, `update sessions s set token_hash = $2, expires_at = $3
		where s.token_hash = $1 and s.expires_at > $4 and s.revoked_at is null
		and exists (select 1 from users u where u.user_id = s.user_id and not u.disabled)
		returning s.user_id`, oldHash, newHash, expires, now).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errInvalidSession
	}
	return userID, err
}

func (r sqlSessionRepository) revoke(ctx context.Context, hash []byte, now time.Time) error {
	result, err := r.db.ExecContext(ctx, `update sessions set revoked_at = $2
		where token_hash = $1 and expires_at > $2 and revoked_at is null`, hash, now)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errInvalidSession
	}
	return nil
}
