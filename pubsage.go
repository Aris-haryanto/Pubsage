package pubsage

import (
	"context"

	"github.com/Aris-haryanto/pubsage/discourse"
)

type pubSage struct {
	messageBroker iMessageBroker
}

func newInit(imb iMessageBroker) *pubSage {
	n := new(pubSage)
	n.messageBroker = imb

	return n
}

// Publish message to stack
func (c *pubSage) Publish(ctx context.Context, message discourse.Publisher) error {
	return c.messageBroker.Publish(ctx, message)
}

// Listen incoming message from stack
func (c *pubSage) Listener(ctx context.Context, fn func(r discourse.Message) error, cfgSubscription discourse.Subscription) error {
	return c.messageBroker.Listener(ctx, fn, cfgSubscription)
}

// Close All connection from the stack
func (c *pubSage) Close() error {
	return c.messageBroker.Close()
}
