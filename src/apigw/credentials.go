package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
	errInvalidUserID      = errors.New("user ID must be 3-64 characters using letters, digits, dots, underscores, or hyphens")
	errInvalidPassword    = errors.New("password must be 12-72 bytes")
	errUserExists         = errors.New("user ID already exists")
	errUserNotFound       = errors.New("user not found")
	errPasswordUnchanged  = errors.New("new password must differ from current password")
	userIDPattern         = regexp.MustCompile(`^[a-z0-9._-]{3,64}$`)
	dummyHash             = []byte("$2a$12$C6UzMDM.H6dfI/f/IKcEeO/2vY.GDln.T7R1Gff7Gy15qVPnWE2ri")
)

const bcryptCost = 12

type credentialRecord struct {
	UserID       string
	PasswordHash string
	Disabled     bool
}

type credentialRepository interface {
	get(context.Context, string) (credentialRecord, error)
	create(context.Context, string, string) error
	setPassword(context.Context, string, string) error
	setDisabled(context.Context, string, bool) error
}

type credentialService struct{ repository credentialRepository }

func canonicalUserID(id string) (string, error) {
	id = strings.ToLower(id)
	if !userIDPattern.MatchString(id) {
		return "", errInvalidUserID
	}
	return id, nil
}

func hashPassword(password string) (string, error) {
	if len(password) < 12 || len(password) > 72 {
		return "", errInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(hash), err
}

func (s credentialService) authenticate(ctx context.Context, id, password string) (string, error) {
	id, err := canonicalUserID(id)
	if err != nil {
		id = ""
	}
	record, lookupErr := s.repository.get(ctx, id)
	if lookupErr != nil && !errors.Is(lookupErr, errUserNotFound) {
		return "", lookupErr
	}
	hash := []byte(record.PasswordHash)
	if lookupErr != nil || err != nil {
		hash = dummyHash
	}
	passwordErr := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if lookupErr != nil || err != nil || record.Disabled || passwordErr != nil {
		return "", errInvalidCredentials
	}
	return record.UserID, nil
}

func (s credentialService) create(ctx context.Context, id, password string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	return s.repository.create(ctx, id, hash)
}

func (s credentialService) changePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	id, err := s.authenticate(ctx, id, oldPassword)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare([]byte(oldPassword), []byte(newPassword)) == 1 {
		return errPasswordUnchanged
	}
	return s.resetPassword(ctx, id, newPassword)
}

func (s credentialService) resetPassword(ctx context.Context, id, newPassword string) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	hash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repository.setPassword(ctx, id, hash)
}

func (s credentialService) setDisabled(ctx context.Context, id string, disabled bool) error {
	id, err := canonicalUserID(id)
	if err != nil {
		return err
	}
	return s.repository.setDisabled(ctx, id, disabled)
}
