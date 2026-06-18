// Package services contains business logic and orchestration
package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/qmessentials/qmessentials/intake/messaging"
	"github.com/qmessentials/qmessentials/intake/models"
	"github.com/qmessentials/qmessentials/intake/repositories"
)

type HashTestResultFunc func(data []byte, key *string) string

type TestResultConsumer interface {
	Subscribe(ctx context.Context, stream, subject string) error
}

type TestResultConsumerNATS struct {
	subscriber     messaging.Subscriber
	repository     repositories.TestResultRepository
	hashKey        *string
	hashTestResult HashTestResultFunc
}

func NewTestResultConsumer(subscriber messaging.Subscriber, repository repositories.TestResultRepository,
	hashKey *string, hashTestResult HashTestResultFunc) *TestResultConsumerNATS {
	return &TestResultConsumerNATS{
		subscriber:     subscriber,
		repository:     repository,
		hashKey:        hashKey,
		hashTestResult: hashTestResult,
	}
}

func (c *TestResultConsumerNATS) Subscribe(ctx context.Context, stream, subject string) error {
	return c.subscriber.Subscribe(ctx, stream, subject, func(ctx context.Context, data []byte) error {
		return c.handleResult(ctx, data)
	})
}

func (c *TestResultConsumerNATS) handleResult(ctx context.Context, data []byte) error {
	var result models.TestResult
	if err := json.Unmarshal(data, &result); err != nil {
		return messaging.NewUnmarshalError(err)
	}
	result.HashValue = HashPayload(data, c.hashKey)
	_, err := c.repository.Add(ctx, &result)
	if err != nil {
		return messaging.NewDatabaseError(err)
	}
	return nil
}

func HashPayload(data []byte, key *string) string {
	if key == nil {
		sum := sha256.Sum256(data)
		return hex.EncodeToString(sum[:])
	}
	mac := hmac.New(sha256.New, []byte(*key))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}
