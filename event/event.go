package event

import (
	"context"
	"time"
)

var DefaultEventBus EventBus

// Event represents a generic event
type Event struct {
	ID          string         `json:"id"`
	Source      string         `json:"source"`
	SpecVersion string         `json:"specversion"`
	Type        string         `json:"type"`
	Data        map[string]any `json:"data"`
	Time        time.Time      `json:"time"`
}

// EventBus defines the interface for publishing and subscribing to events
type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler func(ctx context.Context, event Event)) error
	Close() error // Clean up resources
}
