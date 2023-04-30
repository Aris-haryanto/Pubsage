package redis_adapter

import (
	"github.com/Aris-haryanto/pubsage/discourse"
	"github.com/go-redis/redis"
)

func parseMessageToPubSage(msg *redis.Message) discourse.Message {
	set := discourse.Message{}
	set.Topic = msg.Channel

	// set general object struct
	set.Data = []byte(msg.Payload)

	return set
}
