package nats

import "github.com/nats-io/go-nats"

// Subscriber specifies API for subscribing to NATS topics.
type Subscriber interface {
	// Subscribe subscribes on messages from specific subject within specific queue
	Subscribe(subject string, queue string) (*nats.Subscription, error)
}
