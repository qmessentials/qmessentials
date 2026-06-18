// Package messaging contains messaging-related functionality
package messaging

import (
	"context"
	"errors"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type NatsPublisher struct {
	js jetstream.JetStream
}

func NewNatsPublisher(js jetstream.JetStream) *NatsPublisher {
	return &NatsPublisher{js: js}
}

func (p *NatsPublisher) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := p.js.Publish(ctx, subject, data)
	return err
}

type NatsSubscriber struct {
	js jetstream.JetStream
}

func NewNatsSubscriber(js jetstream.JetStream) *NatsSubscriber {
	return &NatsSubscriber{js: js}
}

func (s *NatsSubscriber) Subscribe(ctx context.Context, stream string, subject string, handler func(ctx context.Context, data []byte) error) error {
	consumer, err := s.js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		Durable:       "workers",
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	iter, err := consumer.Messages()
	if err != nil {
		return err
	}

	go func() {
		defer iter.Stop()
		for {
			msg, err := iter.Next()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				slog.Error("message iterator error", "error", err)
				return
			}
			if err = handler(ctx, msg.Data()); err != nil {
				slog.Error("message handler error", "error", err)
				if errors.Is(err, &UnmarshalError{}) {
					msg.Term()
				} else if errors.Is(err, &DatabaseError{}) {
					msg.Nak()
				} else {
					slog.Error("unknown error", "error", err)
					return
				}
			}
		}
	}()

	go func() {
		<-ctx.Done()
		iter.Stop()
	}()

	return nil
}
