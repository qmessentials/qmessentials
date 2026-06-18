package messaging

import (
	"context"
	"fmt"
)

type UnmarshalError struct {
	err error
}

func NewUnmarshalError(err error) *UnmarshalError {
	return &UnmarshalError{err}
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("failed to unmarshal test result: %s", e.err)
}

type DatabaseError struct {
	err error
}

func NewDatabaseError(err error) *DatabaseError {
	return &DatabaseError{err}
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("failed to save test result: %s", e.err)
}

type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, stream string, subject string, handler func(ctx context.Context, data []byte) error) error
}
