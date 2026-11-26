package messaging

import (
	"context"
	"errors"
)

var UnmarshalError = errors.New("failed to unmarshal test result")
var DatabaseError = errors.New("failed to save test result")

func NewUnmarshalError(err error) error {
	return errors.Join(UnmarshalError, err)
}

func NewDatabaseError(err error) error {
	return errors.Join(DatabaseError, err)
}

type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, stream string, subject string, handler func(ctx context.Context, data []byte) error) error
}
