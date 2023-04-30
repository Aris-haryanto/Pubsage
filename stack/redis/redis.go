package redis_adapter

import (
	"context"
	"fmt"

	"github.com/Aris-haryanto/pubsage/discourse"
	"github.com/go-redis/redis"
)

type RedisAdapter struct {
	conn *redis.Client
}

func (red *RedisAdapter) Close() error {
	if err := red.conn.Close(); err != nil {
		return fmt.Errorf("[Pubsage-Redis] error when close redis connection - %s", err)
	}

	return nil
}

func (red *RedisAdapter) Publish(ctx context.Context, message discourse.Publisher) error {
	pb := red.conn.Publish(message.Topic, message.Data)
	return pb.Err()
}

func (red *RedisAdapter) Listener(ctx context.Context, fn func(discourse.Message) error, cfgSubscription discourse.Subscription) error {
	var (
		getMessage *redis.Message
		errMessage error
	)

	subscriber := red.conn.Subscribe(cfgSubscription.Topic)
	for {
		getMessage, errMessage = subscriber.ReceiveMessage()
		if errMessage != nil {
			break
		}
		message := parseMessageToPubSage(getMessage)

		if errFn := fn(message); errFn != nil {
			fmt.Printf("Redis skip commiting message - Error: %s\n", errFn)
			continue
		}
	}

	return errMessage
}
