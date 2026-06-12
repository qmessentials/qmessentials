package queue

import "context"

type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

type Subscriber interface {
	Subscribe(subject string, handler func(data []byte) error) error
}
