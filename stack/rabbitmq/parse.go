package rabbitmq_adapter

import (
	"github.com/Aris-haryanto/pubsage/discourse"
	amqp "github.com/rabbitmq/amqp091-go"
)

func parsePublisherToRabbitMq(msg discourse.Publisher) amqp.Publishing {
	set := amqp.Publishing{}

	// set general object struct
	set.Body = msg.Data
	set.MessageId = msg.MessageID
	set.Timestamp = msg.PublishTime

	// ==
	set.Headers = msg.Amqp_Headers

	// Properties
	set.ContentType = msg.Amqp_ContentType
	set.ContentEncoding = msg.Amqp_ContentEncoding
	set.DeliveryMode = msg.Amqp_DeliveryMode
	set.Priority = msg.Amqp_Priority
	set.CorrelationId = msg.Amqp_CorrelationId
	set.ReplyTo = msg.Amqp_ReplyTo
	set.Expiration = msg.Amqp_Expiration
	set.Type = msg.Amqp_Type
	set.UserId = msg.Amqp_UserId
	set.AppId = msg.Amqp_AppId

	return set
}

func parseMessageToPubSage(msg amqp.Delivery) discourse.Message {
	set := discourse.Message{}

	// set general object struct
	set.Data = msg.Body
	set.Topic = msg.RoutingKey
	set.MessageID = msg.MessageId
	set.PublishTime = msg.Timestamp

	// ==
	set.Amqp_Headers = msg.Headers

	// Properties
	set.Amqp_ContentType = msg.ContentType
	set.Amqp_ContentEncoding = msg.ContentEncoding
	set.Amqp_DeliveryMode = msg.DeliveryMode
	set.Amqp_Priority = msg.Priority
	set.Amqp_CorrelationId = msg.CorrelationId
	set.Amqp_ReplyTo = msg.ReplyTo
	set.Amqp_Expiration = msg.Expiration
	set.Amqp_Type = msg.Type
	set.Amqp_UserId = msg.UserId
	set.Amqp_AppId = msg.AppId

	// Valid only with Channel.Consume
	set.Amqp_ConsumerTag = msg.ConsumerTag

	// Valid only with Channel.Get
	set.Amqp_MessageCount = msg.MessageCount

	set.Amqp_DeliveryTag = msg.DeliveryTag
	set.Amqp_Redelivered = msg.Redelivered
	set.Amqp_Exchange = msg.Exchange

	return set
}
