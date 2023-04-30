package pubsage

import (
	"context"

	"github.com/Aris-haryanto/pubsage/discourse"
)

type iMessageBroker interface {
	Publish(ctx context.Context, message discourse.Publisher) error
	Listener(ctx context.Context, fn func(r discourse.Message) error, cfgSubscription discourse.Subscription) error
	Close() error
}
