package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserID      = errors.New("user ID must be 3-64 characters using letters, digits, dots, underscores, or hyphens")
	ErrInvalidPassword    = errors.New("password must be 12-72 bytes")
	ErrUserExists         = errors.New("user ID already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordUnchanged  = errors.New("new password must differ from current password")
	userIDPattern         = regexp.MustCompile(`^[a-z0-9._-]{3,64}$`)
	dummyHash             = []byte("$2a$12$C6UzMDM.H6dfI/f/IKcEeO/2vY.GDln.T7R1Gff7Gy15qVPnWE2ri")
)

const bcryptCost = 12

type CredentialRecord struct {
	UserID       string
	PasswordHash string
	Disabled     bool
}

type CredentialRepository interface {
	Get(context.Context, string) (CredentialRecord, error)
	Create(context.Context, string, string) error
	SetPassword(context.Context, string, string) error
	SetDisabled(context.Context, string, bool) error
}

type Credentials struct{ repository CredentialRepository }

func NewCredentials(repository CredentialRepository) Credentials {
	return Credentials{repository: repository}
}

func canonicalUserID(id string) (string, error) {
	id = strings.ToLower(id)
	if !userIDPattern.MatchString(id) {
		return "", ErrInvalidUserID
	}
	return id, nil
}

func hashPassword(password string) (string, error) {
	if len(password) < 12 || len(password) > 72 {
		return "", ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(hash), err
}

func (s Credentials) Authenticate(ctx context.Context, id, password string) (string, error) {
	id, err := canonicalUserID(id)
	if err != nil {
		id = ""
	}
	record, lookupErr := s.repository.Get(ctx, id)
	if lookupErr != nil && !errors.Is(lookupErr, ErrUserNotFound) {
		return "", lookupErr
	}
	hash := []byte(record.PasswordHash)
	if lookupErr != nil || err != nil {
		hash = dummyHash
	}
	passwordErr := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if lookupErr != nil || err != nil || record.Disabled || passwordErr != nil {
		return "", ErrInvalidCredentials
	}
	return record.UserID, nil
}

func (s Credentials) Create(ctx context.Context, id, password string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	return s.repository.Create(ctx, id, hash)
}

func (s Credentials) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	id, err := s.Authenticate(ctx, id, oldPassword)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare([]byte(oldPassword), []byte(newPassword)) == 1 {
		return ErrPasswordUnchanged
	}
	return s.ResetPassword(ctx, id, newPassword)
}

func (s Credentials) ResetPassword(ctx context.Context, id, newPassword string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repository.SetPassword(ctx, id, hash)
}

func (s Credentials) SetDisabled(ctx context.Context, id string, disabled bool) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	return s.repository.SetDisabled(ctx, id, disabled)
}
