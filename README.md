# Pubsage  
[![Go Report Card](https://goreportcard.com/badge/github.com/Aris-haryanto/pubsage)](https://goreportcard.com/report/github.com/Aris-haryanto/pubsage) 
[![GoDoc](https://godoc.org/github.com/Aris-haryanto/pubsage?status.svg)](https://godoc.org/github.com/Aris-haryanto/pubsage)

This is a library adapter to connect your logic code to the stack (pubsub or message broker) with an abstraction.\
so you don't have to change the implementation code when you want to change the stack of pubsub or message broker,
you only need to change the initialization

## Motivation

As we know, there's many pubsub and message broker, and the implementation is also different,
sometime when our code was release to production, somehow we need to change to other stack maybe from Google Pubsub to Kafka, or from RabbitMQ to Redis Pubsub\
and absolutely we need to change our code, and need more time to implement and learn how to implement to new stack

thay's why I created this library, to generalize Producer and Consumer on many pubsub and message broker.

## Stack Support 
- Kafka
- Google Pubsub
- Redis Pubsub
- RabbitMQ

## Installation
just type 
```cli
go get github.com/Aris-haryanto/pubsage
```
Or if you use go module
just import `github.com/Aris-haryanto/pubsage` in your go file
then type
```cli
go mod tidy
```

this will automaticly download the module


## Simple Setup
if you don't want to confuse about setup the connection or channel, you can use this function.\
let Pubsage take the rest of setup until it's ready to use

> **Warning**
> don't forget to call Pubsage Close() function 

### Redis
you need set struct of `redis.Option`, this struct from `github.com/go-redis/redis`
```go
client := pubsage.NewRedis(&redis.Options{
    Addr: "localhost:6379",
})
defer client.Close()
```

### Google Pubsub
you need to set project ID and optionaly set `option.ClientOption`, this struct from library `google.golang.org/api/option`
```go
client := pubsage.NewGooglePubsub(context.Background(), "pubsage-project", option.WithoutAuthentication())
defer client.Close()
```

### Kafka
if you don't want to use Publish/Write function don't fill this param,\
this param is struct from library `github.com/segmentio/kafka-go`
```go
client := pubsage.NewKafka(&kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
	  })
defer client.Close()
```

### Rabbit MQ
set param url with DSN. ie. `amqp://guest:guest@localhost:5672`
```go
client := pubsage.NewRabbitMQ("amqp://guest:guest@localhost:5672")
defer client.Close()
```

## Setup as Adapter

of course you can setup the connection by your self, and use Pubsage as adapter\
call this function and set with your own connection and Pubsage will take the Publish and Subscribe function\
so you can doing other thing with the connection (not only for publish and subscribe)

> **Warning**
> when you use this function, don't call Pubsage Close() function

### Redis
you need to create Redis connection with this library `github.com/go-redis/redis` and put on pubsage.RedisAdapter()
```go
conn := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
client := pubsage.RedisAdapter(conn)
```

### Google Pubsub
you need to create Google Pubsub connection with this library `cloud.google.com/go/pubsub` and put on pubsage.GooglePubsubAdapter()
```go
conn, err := pubsub.NewClient(context.Background(), "pubsage-project", option.WithoutAuthentication())
if err != nil {
    log.Fatalf("%v", err)
}
client := pubsage.GooglePubsubAdapter(conn)
```

### Rabbit MQ
you need to create Rabbit MQ connection with this library `amqp "github.com/rabbitmq/amqp091-go"` and put on pubsage.RabbitMQAdapter()
```go
conn, errConn := amqp.Dial("amqp://guest:guest@localhost:5672")
if errConn != nil {
    log.Fatalf("%v", errConn)
}

ch, errCh := conn.Channel()
if errCh != nil {
    log.Fatalf("%v", errCh)
}
defer conn.Close()

client := pubsage.RabbitMQAdapter(conn)
```

### Kafka
> **Note**
> Kafka doesn't have KafkaAdapter(), like other Because the connection is directly to it's reader or writer


## How To Use

### Publish Message
```go
. . .

attr := map[string][]byte{
    "key_attr": []byte("value attr"),
}

err := client.Publish(context.Background(), discourse.Publisher{
    Topic:      "example-topic",
    Data:       []byte(fmt.Sprintf("this is message - %s", time.Now())),
    Attributes: attr,
})

if err != nil {
    fmt.Println(err)
}
```

### Receive Message
```go
. . .

err := client.Listener(context.Background(), func(msg discourse.Message) error {

            // Your logic here
	    
            fmt.Printf("Message : %s\n", string(msg.Data))
            fmt.Printf("message ID : %s\n", msg.MessageID)
            fmt.Printf("Attribute : %s\n===\n", msg.Attributes)

            // return error then your message will be nack

            // return nil then your message will be ack
            return nil
        }, discourse.Subscription{
            Topic:      "example-topic",
            AutoCommit: false,
        })

if err != nil {
    log.Fatalln(err)
}
```

## About the discourse
discourse is a struct to bridge data from Pubsage to the stack,\
param with Prefix `Amqp_`, `GooglePubsub_`, `Kafka_` this is only use if you initialize connection with one of that stack.\
because this param is feature from the stack it self, can't be the same with other.\
the cool is when you change the initialize to other stack this param will be ignoring, so it safe!

### discourse.Publisher
this struct used to publish the message with some general param

```go
type Publisher struct {
	// which topic message want to consume
	//
	// Available for All
	Topic string

	// Data is the actual data from pubsub/message broker
	//
	// Available for All
	Data []byte

	// PublishTime is the time at which the message was published.
	// on kafka you can set this, but in RabbitMQ and GooglePubsub
	// this is generate from server
	//
	// Available on Kafka
	PublishTime time.Time

	// ID identifies this message. This ID is assigned by the server and is
	// populated for Messages obtained from a subscription.
	//
	// Available for RabbitMQ
	MessageID string

	// Attributes represents the key-value pairs the current message
	// in Kafka this is same as protocol.Header
	// in GooglePubsub this is same as Attributes
	//
	// Available for GooglePubsub, Kafka
	Attributes map[string][]byte

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

	// if not set will generate automatilcy and send to random broker
	Kafka_Key []byte

	// This field is used to hold arbitrary data you wish to include, so it
	// will be available when handle it on the Writer's `Completion` method,
	// this support the application can do any post operation on each message.
	Kafka_WriterData interface{}
}
```

### discourse.Message
this struct used to get the message from the stack with some general param

```go
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
```

### discourse.Subscription
this struct used to configure consumer with some general param

```go
type Subscription struct {
	// The topic to read messages from.
        // 
        // Available for All
	Topic string

	//in Kafka Only used when GroupID is set 
	//
	// Available on Kafka, RabbitMQ
	AutoCommit bool

	// The list of broker addresses used to connect to the kafka cluster.
	Kafka_Brokers []string

	// GroupID holds the optional consumer group id.  If GroupID is specified, then
	// Partition should NOT be specified e.g. 0
	Kafka_GroupID string

	// GroupTopics allows specifying multiple topics, but can only be used in
	// combination with GroupID, as it is a consumer-group feature. As such, if
	// GroupID is set, then either Topic or GroupTopics must be defined.
	Kafka_GroupTopics []string

	// Partition to read messages from.  Either Partition or GroupID may
	// be assigned, but not both
	Kafka_Partition int

	// An dialer used to open connections to the kafka server. This field is
	// optional, if nil, the default dialer is used instead.
	Kafka_Dialer *kafka.Dialer

	// The capacity of the internal message queue, defaults to 100 if none is
	// set.
	Kafka_QueueCapacity int

	// MinBytes indicates to the broker the minimum batch size that the consumer
	// will accept. Setting a high minimum when consuming from a low-volume topic
	// may result in delayed delivery when the broker does not have enough data to
	// satisfy the defined minimum.
	//
	// Default: 1
	Kafka_MinBytes int

	// MaxBytes indicates to the broker the maximum batch size that the consumer
	// will accept. The broker will truncate a message to satisfy this maximum, so
	// choose a value that is high enough for your largest message size.
	//
	// Default: 1MB
	Kafka_MaxBytes int

	// Maximum amount of time to wait for new data to come when fetching batches
	// of messages from kafka.
	//
	// Default: 10s
	Kafka_MaxWait time.Duration

	// ReadBatchTimeout amount of time to wait to fetch message from kafka messages batch.
	//
	// Default: 10s
	Kafka_ReadBatchTimeout time.Duration

	// ReadLagInterval sets the frequency at which the reader lag is updated.
	// Setting this field to a negative value disables lag reporting.
	Kafka_ReadLagInterval time.Duration

	// GroupBalancers is the priority-ordered list of client-side consumer group
	// balancing strategies that will be offered to the coordinator.  The first
	// strategy that all group members support will be chosen by the leader.
	//
	// Default: [Range, RoundRobin]
	//
	// Only used when GroupID is set
	Kafka_GroupBalancers []kafka.GroupBalancer

	// HeartbeatInterval sets the optional frequency at which the reader sends the consumer
	// group heartbeat update.
	//
	// Default: 3s
	//
	// Only used when GroupID is set
	Kafka_HeartbeatInterval time.Duration

	// PartitionWatchInterval indicates how often a reader checks for partition changes.
	// If a reader sees a partition change (such as a partition add) it will rebalance the group
	// picking up new partitions.
	//
	// Default: 5s
	//
	// Only used when GroupID is set and WatchPartitionChanges is set.
	Kafka_PartitionWatchInterval time.Duration

	// WatchForPartitionChanges is used to inform kafka-go that a consumer group should be
	// polling the brokers and rebalancing if any partition changes happen to the topic.
	Kafka_WatchPartitionChanges bool

	// SessionTimeout optionally sets the length of time that may pass without a heartbeat
	// before the coordinator considers the consumer dead and initiates a rebalance.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	Kafka_SessionTimeout time.Duration

	// RebalanceTimeout optionally sets the length of time the coordinator will wait
	// for members to join as part of a rebalance.  For kafka servers under higher
	// load, it may be useful to set this value higher.
	//
	// Default: 30s
	//
	// Only used when GroupID is set
	Kafka_RebalanceTimeout time.Duration

	// JoinGroupBackoff optionally sets the length of time to wait between re-joining
	// the consumer group after an error.
	//
	// Default: 5s
	Kafka_JoinGroupBackoff time.Duration

	// RetentionTime optionally sets the length of time the consumer group will be saved
	// by the broker
	//
	// Default: 24h
	//
	// Only used when GroupID is set
	Kafka_RetentionTime time.Duration

	// StartOffset determines from whence the consumer group should begin
	// consuming when it finds a partition without a committed offset.  If
	// non-zero, it must be set to one of FirstOffset or LastOffset.
	//
	// Default: FirstOffset
	//
	// Only used when GroupID is set
	Kafka_StartOffset int64

	// BackoffDelayMin optionally sets the smallest amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 100ms
	Kafka_ReadBackoffMin time.Duration

	// BackoffDelayMax optionally sets the maximum amount of time the reader will wait before
	// polling for new messages
	//
	// Default: 1s
	Kafka_ReadBackoffMax time.Duration

	// If not nil, specifies a logger used to report internal changes within the
	// reader.
	Kafka_Logger kafka.Logger

	// ErrorLogger is the logger used to report errors. If nil, the reader falls
	// back to using Logger instead.
	Kafka_ErrorLogger kafka.Logger

	// IsolationLevel controls the visibility of transactional records.
	// ReadUncommitted makes all records visible. With ReadCommitted only
	// non-transactional and committed records are visible.
	Kafka_IsolationLevel kafka.IsolationLevel

	// Limit of how many attempts to connect will be made before returning the error.
	//
	// The default is to try 3 times.
	Kafka_MaxAttempts int

	// OffsetOutOfRangeError indicates that the reader should return an error in
	// the event of an OffsetOutOfRange error, rather than retrying indefinitely.
	// This flag is being added to retain backwards-compatibility, so it will be
	// removed in a future version of kafka-go.
	Kafka_OffsetOutOfRangeError bool
}
```

# Community 
Feel free to contribute on this library and do not hesitate to open an issue if you want to discuss about a feature.

