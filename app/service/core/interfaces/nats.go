package interfaces

import "github.com/nats-io/nats.go"

type NatsService interface {
	Close()
	Reload(endpoints string) error

	SubscribeQueue(subject, queue string, handler nats.Handler) error
	Subscribe(subject string, handler nats.Handler) error
	Publish(subject string, v interface{}) error
}
