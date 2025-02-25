package messaging

// SubscriptionOptions holds configuration for subscription.
type SubscriptionOptions struct {
	Group string
	Name  string
}

// SubscriptionOption defines a function to set subscription options.
type SubscriptionOption func(opts *SubscriptionOptions)

// WithConsumerGroup specifies the consumer group.
func WithConsumerGroup(group string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Group = group
	}
}

// WithConsumerName specifies the consumer name.
func WithConsumerName(name string) SubscriptionOption {
	return func(opts *SubscriptionOptions) {
		opts.Name = name
	}
}
