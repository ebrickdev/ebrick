package messaging

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

var DefaultEventBus EventBus

// EventBus defines the interface for publishing and subscribing to events
type EventBus interface {
	Publish(ctx context.Context, topic string, event cloudevents.Event) error
	Subscribe(topic string, handler func(ctx context.Context, event cloudevents.Event), opts ...SubscriptionOption) error
	Close() error // Clean up resources
}
