package googlepubsub_adapter

import (
	"context"
	"fmt"

	"github.com/Aris-haryanto/pubsage/discourse"

	"cloud.google.com/go/pubsub"
)

type GooglePubsubAdapter struct {
	conn *pubsub.Client
}

func (gpb *GooglePubsubAdapter) Close() error {
	if err := gpb.conn.Close(); err != nil {
		return fmt.Errorf("[Pubsage-GooglePubsub] error when close google pubsub connection - %s", err)
	}

	return nil
}

func (gpb *GooglePubsubAdapter) Publish(ctx context.Context, message discourse.Publisher) error {
	publisher := gpb.conn.Topic(message.Topic)

	msg := parsePublisherToGooglePubsub(message)

	pb := publisher.Publish(ctx, msg)
	_, err := pb.Get(ctx)

	return err
}

func (gpb *GooglePubsubAdapter) Listener(ctx context.Context, fn func(discourse.Message) error, cfgSubcription discourse.Subscription) error {
	subscriber := gpb.conn.Subscription(cfgSubcription.Topic)

	err := subscriber.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {

		message := parseToMessagePubSage(msg)

		if errFn := fn(message); errFn != nil {
			fmt.Printf("GooglePubsub skip commiting message with ID %s - Error: %s\n", message.MessageID, errFn)
			// Nack when error
			msg.Nack()
		} else {
			// Ack if no error
			msg.Ack()
		}

	})

	return err
}
