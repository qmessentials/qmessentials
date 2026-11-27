package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type SQLUserRepository struct{ db *sql.DB }

func NewSQLUserRepository(db *sql.DB) UserRepository {
	return SQLUserRepository{db: db}
}

func (r SQLUserRepository) List(ctx context.Context) ([]User, error) {
	rows, err := r.db.QueryContext(ctx, `select u.user_id, u.disabled, coalesce(string_agg(ur.role, ',' order by ur.role), '')
		from users u left join user_roles ur on ur.user_id = u.user_id group by u.user_id order by u.user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		var u User
		var roles string
		if err := rows.Scan(&u.UserID, &u.Disabled, &roles); err != nil {
			return nil, err
		}
		u.Roles = splitRoles(roles)
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r SQLUserRepository) Get(ctx context.Context, id string) (User, error) {
	var u User
	var roles string
	err := r.db.QueryRowContext(ctx, `select u.user_id, u.disabled, coalesce(string_agg(ur.role, ',' order by ur.role), '')
		from users u left join user_roles ur on ur.user_id = u.user_id where u.user_id = $1 group by u.user_id`, id).
		Scan(&u.UserID, &u.Disabled, &roles)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	u.Roles = splitRoles(roles)
	return u, err
}

func splitRoles(value string) []string {
	if value == "" {
		return []string{}
	}
	return strings.Split(value, ",")
}

func insertRoles(ctx context.Context, tx *sql.Tx, id string, roles []string) error {
	for _, role := range roles {
		if _, err := tx.ExecContext(ctx, `insert into user_roles (user_id, role) values ($1, $2)`, id, role); err != nil {
			return err
		}
	}
	return nil
}

func (r SQLUserRepository) Create(ctx context.Context, id, hash string, roles []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `insert into users (user_id, password_hash) values ($1, $2)`, id, hash)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrUserExists
	}
	if err != nil {
		return err
	}
	if err := insertRoles(ctx, tx, id, roles); err != nil {
		return err
	}
	return tx.Commit()
}

func lockUserManagement(ctx context.Context, tx *sql.Tx) error {
	// Serialize bootstrap and role changes so two requests cannot each remove
	// a different Administrator after both observe another one remaining.
	_, err := tx.ExecContext(ctx, `select pg_advisory_xact_lock(325)`)
	return err
}

func lastAdministrator(ctx context.Context, tx *sql.Tx, id string) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx, `select count(*) from user_roles ur join users u on u.user_id = ur.user_id
		where ur.role = $2 and not u.disabled and u.user_id <> $1`, id, AdministratorRole).Scan(&count)
	return count == 0, err
}

func (r SQLUserRepository) SetDisabled(ctx context.Context, id string, disabled bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := lockUserManagement(ctx, tx); err != nil {
		return err
	}
	var hasAdmin bool
	err = tx.QueryRowContext(ctx, `select exists(select 1 from user_roles where user_id = $1 and role = $2) from users where user_id = $1`, id, AdministratorRole).Scan(&hasAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if disabled && hasAdmin {
		last, err := lastAdministrator(ctx, tx, id)
		if err != nil {
			return err
		}
		if last {
			return ErrLastAdministrator
		}
	}
	_, err = tx.ExecContext(ctx, `update users set disabled = $2, updated_at = now() where user_id = $1`, id, disabled)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r SQLUserRepository) SetRoles(ctx context.Context, id string, roles []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := lockUserManagement(ctx, tx); err != nil {
		return err
	}
	var disabled bool
	err = tx.QueryRowContext(ctx, `select disabled from users where user_id = $1`, id).Scan(&disabled)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}
	if !disabled {
		var hadAdmin bool
		if err := tx.QueryRowContext(ctx, `select exists(select 1 from user_roles where user_id = $1 and role = $2)`, id, AdministratorRole).Scan(&hadAdmin); err != nil {
			return err
		}
		if hadAdmin && !containsRole(roles, AdministratorRole) {
			last, err := lastAdministrator(ctx, tx, id)
			if err != nil {
				return err
			}
			if last {
				return ErrLastAdministrator
			}
		}
	}
	if _, err := tx.ExecContext(ctx, `delete from user_roles where user_id = $1`, id); err != nil {
		return err
	}
	if err := insertRoles(ctx, tx, id, roles); err != nil {
		return err
	}
	return tx.Commit()
}

func containsRole(roles []string, sought string) bool {
	for _, role := range roles {
		if role == sought {
			return true
		}
	}
	return false
}

func (r SQLUserRepository) Bootstrap(ctx context.Context, id, hash string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := lockUserManagement(ctx, tx); err != nil {
		return err
	}
	var count int
	if err := tx.QueryRowContext(ctx, `select count(*) from user_roles where role = $1`, AdministratorRole).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `insert into users (user_id, password_hash) values ($1, $2)`, id, hash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserExists
		}
		return err
	}
	if err := insertRoles(ctx, tx, id, []string{AdministratorRole}); err != nil {
		return err
	}
	return tx.Commit()
}
