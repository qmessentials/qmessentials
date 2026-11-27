package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"
)

const sessionLifetime = 12 * time.Hour
const SessionCookieName = "session"

var ErrInvalidSession = errors.New("invalid session")

type SessionRepository interface {
	Create(context.Context, []byte, string, time.Time) error
	Get(context.Context, []byte, time.Time) (string, error)
	Renew(context.Context, []byte, []byte, time.Time, time.Time) (string, error)
	Revoke(context.Context, []byte, time.Time) error
}

type Sessions struct {
	repository SessionRepository
	now        func() time.Time
}

func NewSessions(repository SessionRepository) Sessions {
	return Sessions{repository: repository}
}

func (s Sessions) WithClock(now func() time.Time) Sessions {
	s.now = now
	return s
}

func (s Sessions) clock() time.Time {
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
		return nil, ErrInvalidSession
	}
	hash := sha256.Sum256(bytes)
	return hash[:], nil
}

func (s Sessions) Create(ctx context.Context, userID string) (string, time.Time, error) {
	token, hash, err := newSessionToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expires := s.clock().Add(sessionLifetime)
	if err := s.repository.Create(ctx, hash, userID, expires); err != nil {
		return "", time.Time{}, err
	}
	return token, expires, nil
}

func (s Sessions) Validate(ctx context.Context, token string) (string, error) {
	hash, err := hashSessionToken(token)
	if err != nil {
		return "", err
	}
	return s.repository.Get(ctx, hash, s.clock())
}

func (s Sessions) Renew(ctx context.Context, token string) (string, string, time.Time, error) {
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
	userID, err := s.repository.Renew(ctx, oldHash, newHash, now, expires)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return userID, newToken, expires, nil
}

func (s Sessions) Revoke(ctx context.Context, token string) error {
	hash, err := hashSessionToken(token)
	if err != nil {
		return err
	}
	return s.repository.Revoke(ctx, hash, s.clock())
}
