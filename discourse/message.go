package discourse

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	// this message from which Topic
	//
	// Available for Kafka, RabbitMQ, Redis
	Topic string

	// Data is the actual data from pubsub/message broker
	//
	// Available for All
	Data []byte

	// PublishTime is the time at which the message was published.
	// on kafka you can set this, but in RabbitMQ and GooglePubsub
	// this is generate from server
	//
	// Available on Kafka, RabbitMQ, GooglePubsub
	PublishTime time.Time

	// ID identifies this message. This ID is assigned by the server and is
	// populated for Messages obtained from a subscription.
	//
	// Available for Kafka, RabbitMQ, GooglePubsub
	MessageID string

	// Attributes represents the key-value pairs the current message
	// in Kafka this is same as protocol.Header
	// in GooglePubsub this is same as Attributes
	//
	// Available for GooglePubsub, Kafka
	Attributes map[string][]byte

	// DeliveryAttempt is the number of times a message has been delivered.
	// This is part of the dead lettering feature that forwards messages that
	// fail to be processed (from nack/ack deadline timeout) to a dead letter topic.
	// If dead lettering is enabled, this will be set on all attempts, starting
	// with value 1. Otherwise, the value will be nil.
	GooglePubsub_DeliveryAttempt *int

	// OrderingKey identifies related messages for which publish order should
	// be respected. If empty string is used, message will be sent unordered.
	GooglePubsub_OrderingKey string

	Amqp_Headers amqp.Table // Application or header exchange table

	// Properties
	Amqp_ContentType     string // MIME content type
	Amqp_ContentEncoding string // MIME content encoding
	Amqp_DeliveryMode    uint8  // queue implementation use - non-persistent (1) or persistent (2)
	Amqp_Priority        uint8  // queue implementation use - 0 to 9
	Amqp_CorrelationId   string // application use - correlation identifier
	Amqp_ReplyTo         string // application use - address to reply to (ex: RPC)
	Amqp_Expiration      string // implementation use - message expiration spec
	Amqp_Type            string // application use - message type name
	Amqp_UserId          string // application use - creating user - should be authenticated user
	Amqp_AppId           string // application use - creating application id

	// Valid only with Channel.Consume
	Amqp_ConsumerTag string

	// Valid only with Channel.Get
	Amqp_MessageCount uint32

	Amqp_DeliveryTag uint64
	Amqp_Redelivered bool
	Amqp_Exchange    string // basic.publish exchange

	Kafka_Partition     int
	Kafka_Offset        int64
	Kafka_HighWaterMark int64
	Kafka_Key           []byte

	// This field is used to hold arbitrary data you wish to include, so it
	// will be available when handle it on the Writer's `Completion` method,
	// this support the application can do any post operation on each message.
	Kafka_WriterData interface{}
}
