package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// EventWithCtx couples an event with its associated context.
type EventWithCtx struct {
	ctx   context.Context
	event Event
}

// subscriber represents a single event handler.
type subscriber struct {
	id      string
	channel chan EventWithCtx
}

// MemoryEventBus is an in-memory implementation of the EventBus using channels.
type MemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]subscriber
	closed      bool
}

// NewEventBus creates a new MemoryEventBus.
func NewMemoryEventBus() (*MemoryEventBus, error) {
	return &MemoryEventBus{
		subscribers: make(map[string][]subscriber),
	}, nil
}

// generateUniqueID generates a unique identifier for a subscriber.
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Publish sends an event to all subscribers of the specified event type asynchronously.
func (b *MemoryEventBus) Publish(ctx context.Context, topic string, event Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return errors.New("eventbus is closed")
	}

	if topic == "" || event.ID == "" {
		return errors.New("event must have a valid ID and Topic")
	}

	if chans, exists := b.subscribers[topic]; exists {
		for _, sub := range chans {
			go func(c chan EventWithCtx) {
				select {
				case c <- EventWithCtx{ctx: ctx, event: event}:
				case <-ctx.Done():
				}
			}(sub.channel)
		}
	}
	return nil
}

// Subscribe registers a handler for the specified event type.
// It returns an error if the bus is closed.
func (b *MemoryEventBus) Subscribe(topic string, handler func(ctx context.Context, event Event), opts ...SubscriptionOption) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return errors.New("eventbus is closed")
	}

	options := &SubscriptionOptions{}
	for _, opt := range opts {
		opt(options)
	}

	ch := make(chan EventWithCtx, 10) // Buffered channel to prevent blocking
	id := generateUniqueID()
	sub := subscriber{id: id, channel: ch}
	b.subscribers[topic] = append(b.subscribers[topic], sub)

	go func() {
		for e := range ch {
			handler(e.ctx, e.event)
		}
	}()

	return nil
}

// Close shuts down the event bus and cleans up all channels.
func (b *MemoryEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return errors.New("eventbus is already closed")
	}

	b.closed = true
	for _, subs := range b.subscribers {
		for _, sub := range subs {
			close(sub.channel)
		}
	}
	// Clear the subscribers map
	b.subscribers = make(map[string][]subscriber)
	return nil
}
