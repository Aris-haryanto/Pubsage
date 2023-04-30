package googlepubsub_adapter

import (
	"cloud.google.com/go/pubsub"
	"github.com/Aris-haryanto/pubsage/discourse"
)

func parsePublisherToGooglePubsub(message discourse.Publisher) *pubsub.Message {
	set := &pubsub.Message{}

	set.Data = message.Data
	set.OrderingKey = message.GooglePubsub_OrderingKey

	attr := map[string]string{}
	for k, v := range message.Attributes {
		attr[k] = string(v)
	}

	set.Attributes = attr

	return set
}

func parseToMessagePubSage(msg *pubsub.Message) discourse.Message {
	set := discourse.Message{}
	// set general object struct
	set.Data = msg.Data
	set.MessageID = msg.ID
	set.PublishTime = msg.PublishTime

	attr := map[string][]byte{}
	for k, v := range msg.Attributes {
		attr[k] = []byte(v)
	}
	set.Attributes = attr

	// ==
	set.GooglePubsub_DeliveryAttempt = msg.DeliveryAttempt
	set.GooglePubsub_OrderingKey = msg.OrderingKey

	return set
}
