package googlepubsub_adapter

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func New(ctx context.Context, projectID string, opts ...option.ClientOption) *GooglePubsubAdapter {
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		log.Fatalf("[Pubsage-GooglePubsub] %v", err)
	}

	topic, _ := client.CreateTopic(ctx, "kocak")
	client.CreateSubscription(ctx, "kocak-sub", pubsub.SubscriptionConfig{Topic: topic})

	return &GooglePubsubAdapter{conn: client}
}

func NewConn(c *pubsub.Client) *GooglePubsubAdapter {
	return &GooglePubsubAdapter{conn: c}
}
