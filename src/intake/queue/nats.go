package queue

import (
	"context"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(nc *nats.Conn) *NatsPublisher {
	return &NatsPublisher{nc: nc}
}

func (p *NatsPublisher) Publish(ctx context.Context, subject string, data []byte) error {
	return p.nc.Publish(subject, data)
}

type NatsSubscriber struct {
	nc *nats.Conn
}

func NewNatsSubscriber(nc *nats.Conn) *NatsSubscriber {
	return &NatsSubscriber{nc: nc}
}

func (s *NatsSubscriber) Subscribe(subject string, handler func(data []byte) error) error {
	_, err := s.nc.Subscribe(subject, func(m *nats.Msg) {
		if err := handler(m.Data); err != nil {
			// Errors should be handled by the handler or logged here if needed.
			// The interface doesn't specify how to return errors from the callback.
		}
	})
	return err
}
