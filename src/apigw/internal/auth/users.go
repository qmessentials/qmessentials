package auth

import (
	"context"
	"errors"
	"slices"
)

const (
	AdministratorRole  = "Administrator"
	QualityManagerRole = "Quality Manager"
	TechnicianRole     = "Technician"
)

var (
	ErrInvalidRoles      = errors.New("roles must be Administrator, Quality Manager, or Technician")
	ErrLastAdministrator = errors.New("cannot remove or disable the last Administrator")
)

type User struct {
	UserID   string   `json:"userId"`
	Disabled bool     `json:"disabled"`
	Roles    []string `json:"roles"`
}

type UserRepository interface {
	List(context.Context) ([]User, error)
	Get(context.Context, string) (User, error)
	Create(context.Context, string, string, []string) error
	SetDisabled(context.Context, string, bool) error
	SetRoles(context.Context, string, []string) error
	Bootstrap(context.Context, string, string) error
}

type Users struct{ repository UserRepository }

func NewUsers(repository UserRepository) Users {
	return Users{repository: repository}
}

func (s Users) List(ctx context.Context) ([]User, error) {
	return s.repository.List(ctx)
}

func (s Users) Get(ctx context.Context, id string) (User, error) {
	return s.repository.Get(ctx, id)
}

func validateRoles(roles []string) error {
	if roles == nil {
		return ErrInvalidRoles
	}
	seen := map[string]bool{}
	for _, role := range roles {
		switch role {
		case AdministratorRole, QualityManagerRole, TechnicianRole:
		default:
			return ErrInvalidRoles
		}
		if seen[role] {
			return ErrInvalidRoles
		}
		seen[role] = true
	}
	return nil
}

func (s Users) Create(ctx context.Context, id, password string, roles []string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	if err := validateRoles(roles); err != nil {
		return err
	}
	return s.repository.Create(ctx, id, hash, roles)
}

func (s Users) SetRoles(ctx context.Context, id string, roles []string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	if err := validateRoles(roles); err != nil {
		return err
	}
	return s.repository.SetRoles(ctx, id, roles)
}

func (s Users) SetDisabled(ctx context.Context, id string, disabled bool) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	return s.repository.SetDisabled(ctx, id, disabled)
}

func (s Users) Bootstrap(ctx context.Context, id, password string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	return s.repository.Bootstrap(ctx, id, hash)
}

func (s Users) IsAdministrator(ctx context.Context, id string) (bool, error) {
	u, err := s.repository.Get(ctx, id)
	if err != nil {
		return false, err
	}
	return !u.Disabled && slices.Contains(u.Roles, AdministratorRole), nil
}
