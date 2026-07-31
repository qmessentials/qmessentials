// Package repositories provides database access.
package repositories

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qmessentials/qmessentials/subscription/models"
)

type SubscriptionRepository interface {
	GetForUser(ctx context.Context, ownerUserID string, activeOnly bool) ([]models.Subscription, error)
	GetActive(ctx context.Context) ([]models.Subscription, error)
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	Create(ctx context.Context, ownerUserID string, input models.CreateSubscriptionInput) (*models.Subscription, error)
	Update(ctx context.Context, id int, input models.UpdateSubscriptionInput) (*models.Subscription, error)
	Deactivate(ctx context.Context, id int) error
}

func (r *SubscriptionRepositoryPG) GetActive(ctx context.Context) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx, `
		select id, owner_user_id, rule_text, version_id, is_active, created_at, updated_at
		from subscriptions
		where is_active
		order by id`)
	if err != nil {
		slog.Error("failed to query active subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	results := make([]models.Subscription, 0)
	for rows.Next() {
		var subscription models.Subscription
		if err = rows.Scan(
			&subscription.ID,
			&subscription.OwnerUserID,
			&subscription.RuleText,
			&subscription.VersionID,
			&subscription.IsActive,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			slog.Error("failed to scan active subscription row", "error", err)
			return nil, err
		}
		results = append(results, subscription)
	}
	return results, rows.Err()
}

func (r *SubscriptionRepositoryPG) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		select id, owner_user_id, rule_text, version_id, is_active, created_at, updated_at
		from subscriptions
		where id = $1`, id)

	var subscription models.Subscription
	if err := row.Scan(
		&subscription.ID,
		&subscription.OwnerUserID,
		&subscription.RuleText,
		&subscription.VersionID,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		slog.Error("failed to get subscription", "error", err)
		return nil, err
	}
	return &subscription, nil
}

var ErrSubscriptionNotFound = errors.New("subscription not found")

func (r *SubscriptionRepositoryPG) Create(ctx context.Context, ownerUserID string, input models.CreateSubscriptionInput) (*models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		insert into subscriptions (owner_user_id, rule_text, version_id, is_active)
		values ($1, $2, 1, true)
		returning id, owner_user_id, rule_text, version_id, is_active, created_at, updated_at`,
		ownerUserID,
		input.RuleText,
	)

	var subscription models.Subscription
	if err := row.Scan(
		&subscription.ID,
		&subscription.OwnerUserID,
		&subscription.RuleText,
		&subscription.VersionID,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	); err != nil {
		slog.Error("failed to create subscription", "error", err)
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepositoryPG) Update(ctx context.Context, id int, input models.UpdateSubscriptionInput) (*models.Subscription, error) {
	row := r.db.QueryRow(ctx, `
		update subscriptions
		set rule_text = $1,
		    version_id = version_id + 1,
		    updated_at = now()
		where id = $2
		returning id, owner_user_id, rule_text, version_id, is_active, created_at, updated_at`,
		input.RuleText,
		id,
	)

	var subscription models.Subscription
	if err := row.Scan(
		&subscription.ID,
		&subscription.OwnerUserID,
		&subscription.RuleText,
		&subscription.VersionID,
		&subscription.IsActive,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		slog.Error("failed to update subscription", "error", err)
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepositoryPG) Deactivate(ctx context.Context, id int) error {
	commandTag, err := r.db.Exec(ctx, `
		update subscriptions
		set is_active = false,
		    updated_at = now()
		where id = $1`,
		id,
	)
	if err != nil {
		slog.Error("failed to deactivate subscription", "error", err)
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}
	return nil
}

type SubscriptionRepositoryPG struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepositoryPG(db *pgxpool.Pool) *SubscriptionRepositoryPG {
	return &SubscriptionRepositoryPG{db: db}
}

func (r *SubscriptionRepositoryPG) GetForUser(ctx context.Context, ownerUserID string, activeOnly bool) ([]models.Subscription, error) {
	rows, err := r.db.Query(ctx, `
		select id, owner_user_id, rule_text, version_id, is_active, created_at, updated_at
		from subscriptions
		where owner_user_id = $1
		  and (not $2 or is_active)
		order by id`, ownerUserID, activeOnly)
	if err != nil {
		slog.Error("failed to query subscriptions", "error", err)
		return nil, err
	}
	defer rows.Close()

	results := make([]models.Subscription, 0)
	for rows.Next() {
		var subscription models.Subscription
		if err = rows.Scan(
			&subscription.ID,
			&subscription.OwnerUserID,
			&subscription.RuleText,
			&subscription.VersionID,
			&subscription.IsActive,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			slog.Error("failed to scan subscription row", "error", err)
			return nil, err
		}
		results = append(results, subscription)
	}
	return results, rows.Err()
}
