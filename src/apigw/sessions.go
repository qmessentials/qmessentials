package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"
)

const sessionLifetime = 12 * time.Hour
const sessionCookieName = "session"

var errInvalidSession = errors.New("invalid session")

type sessionRepository interface {
	create(context.Context, []byte, string, time.Time) error
	get(context.Context, []byte, time.Time) (string, error)
	renew(context.Context, []byte, []byte, time.Time, time.Time) (string, error)
	revoke(context.Context, []byte, time.Time) error
}

type sessionService struct {
	repository sessionRepository
	now        func() time.Time
}

func (s sessionService) clock() time.Time {
	if s.now != nil {
		return s.now().UTC()
	}
	return time.Now().UTC()
}

func newSessionToken() (string, []byte, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	hash := sha256.Sum256(bytes)
	return token, hash[:], nil
}

func hashSessionToken(token string) ([]byte, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(bytes) != 32 || base64.RawURLEncoding.EncodeToString(bytes) != token {
		return nil, errInvalidSession
	}
	hash := sha256.Sum256(bytes)
	return hash[:], nil
}

func (s sessionService) create(ctx context.Context, userID string) (string, time.Time, error) {
	token, hash, err := newSessionToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expires := s.clock().Add(sessionLifetime)
	if err := s.repository.create(ctx, hash, userID, expires); err != nil {
		return "", time.Time{}, err
	}
	return token, expires, nil
}

func (s sessionService) validate(ctx context.Context, token string) (string, error) {
	hash, err := hashSessionToken(token)
	if err != nil {
		return "", err
	}
	return s.repository.get(ctx, hash, s.clock())
}

func (s sessionService) renew(ctx context.Context, token string) (string, string, time.Time, error) {
	oldHash, err := hashSessionToken(token)
	if err != nil {
		return "", "", time.Time{}, err
	}
	newToken, newHash, err := newSessionToken()
	if err != nil {
		return "", "", time.Time{}, err
	}
	now := s.clock()
	expires := now.Add(sessionLifetime)
	userID, err := s.repository.renew(ctx, oldHash, newHash, now, expires)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return userID, newToken, expires, nil
}

func (s sessionService) revoke(ctx context.Context, token string) error {
	hash, err := hashSessionToken(token)
	if err != nil {
		return err
	}
	return s.repository.revoke(ctx, hash, s.clock())
}
