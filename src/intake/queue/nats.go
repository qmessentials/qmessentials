// Package queue contains queue-related functionality
package queue

import (
	"context"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	js nats.JetStreamContext
}

func NewNatsPublisher(js nats.JetStreamContext) *NatsPublisher {
	return &NatsPublisher{js: js}
}

func (p *NatsPublisher) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := p.js.Publish(subject, data)
	return err
}

type NatsSubscriber struct {
	js nats.JetStreamContext
}

func NewNatsSubscriber(js nats.JetStreamContext) *NatsSubscriber {
	return &NatsSubscriber{js: js}
}

func (s *NatsSubscriber) Subscribe(subject string, handler func(data []byte) error) error {
	_, err := s.js.QueueSubscribe(subject, "workers", func(m *nats.Msg) {
		if err := handler(m.Data); err != nil {
			_ = m.Nak()
		} else {
			_ = m.Ack()
		}
	}, nats.Durable("workers"))
	return err
}
