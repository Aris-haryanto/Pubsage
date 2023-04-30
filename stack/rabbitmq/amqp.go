package rabbitmq_adapter

import (
	"context"
	"fmt"

	"github.com/Aris-haryanto/pubsage/discourse"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (rmq *RabbitMQAdapter) Close() error {
	if err := rmq.conn.Close(); err != nil {
		return fmt.Errorf("[Pubsage-RabbitMQ] error when close rabbitmq connection - %s", err)
	}
	if err := rmq.ch.Close(); err != nil {
		return fmt.Errorf("[Pubsage-RabbitMQ] error when close channel rabbitmq - %s", err)
	}

	return nil
}

func (rmq *RabbitMQAdapter) Publish(ctx context.Context, message discourse.Publisher) error {
	msg := parsePublisherToRabbitMq(message)
	err := rmq.ch.PublishWithContext(ctx, "", message.Topic, false, false, msg)

	return err
}

func (rmq *RabbitMQAdapter) Listener(ctx context.Context, fn func(discourse.Message) error, cfgSubcription discourse.Subscription) error {
	_, errQueue := rmq.ch.QueueDeclare(
		cfgSubcription.Topic, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if errQueue != nil {
		return errQueue
	}

	messages, errConsume := rmq.ch.Consume(
		cfgSubcription.Topic,      // queue
		"",                        // consumer
		cfgSubcription.AutoCommit, // auto-ack
		false,                     // exclusive
		false,                     // no-local
		false,                     // no-wait
		nil,                       // args
	)
	if errConsume != nil {
		return errConsume
	}

	go func() {
		for msg := range messages {
			parse := parseMessageToPubSage(msg)
			if errFn := fn(parse); errFn != nil {
				fmt.Printf("[Pubsage-RabbitMQ] skip commiting message with key %s - Error: %s\n", parse.MessageID, errFn)
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()

	var forever chan struct{}
	<-forever

	return nil
}
