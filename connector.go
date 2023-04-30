package pubsage

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/go-redis/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/api/option"

	googlepubsub_adapter "github.com/Aris-haryanto/pubsage/stack/googlepubsub"
	kafka_adapter "github.com/Aris-haryanto/pubsage/stack/kafka"
	rabbitmq_adapter "github.com/Aris-haryanto/pubsage/stack/rabbitmq"
	redis_adapter "github.com/Aris-haryanto/pubsage/stack/redis"
)

// If you have a publisher, you have to fill kafka.Write
// If not you can leave it
//
// don't forget to call Close() in defer, when you use this function with filled param
//
// NOTE: Kafka doesn't have KafkaAdapter(), like other Because the connection is directly to it's reader or writer
func NewKafka(w *kafka.Writer) *pubSage {
	conn := kafka_adapter.New(w)
	return newInit(conn)
}

// if you don't want to confuse setting up the connection and channel
// you can use this function,
// just set param url with DSN. ie. amqp://guest:guest@localhost:5672
// then Pubsage will take the rest of configuration until it's ready to use
//
// don't forget to call Close() in defer, when you use this function
func NewRabbitMQ(url string) *pubSage {
	conn := rabbitmq_adapter.New(url)
	return newInit(conn)
}

// you can setup the connection and channel by your self
// and use Pubsage as adapter
// call this function and set with your amqp channel
// and Pubsage will take the Publish and Listener
//
// when you use this function , don't call Close()
func RabbitMQAdapter(ch *amqp.Channel) *pubSage {
	conn := rabbitmq_adapter.NewConn(ch)
	return newInit(conn)
}

// if you don't want to confuse setting up the connection
// you can use this function,
// just set param url with struct redis.Options
// then Pubsage will take the rest of configuration until it's ready to use
//
// don't forget to call Close() in defer, when you use this function
func NewRedis(c *redis.Options) *pubSage {
	conn := redis_adapter.New(c)
	return newInit(conn)
}

// you can setup the connection by your self
// and use Pubsage as adapter
// call this function and set with your redis connection
// and Pubsage will take the Publish and Listener
//
// when you use this function, don't call Close()
func RedisAdapter(c *redis.Client) *pubSage {
	conn := redis_adapter.NewConn(c)
	return newInit(conn)
}

// if you don't want to confuse setting up the connection
// you can use this function,
// just set param context, project_id and option
// then Pubsage will take the rest of configuration until it's ready to use
//
// don't forget to call Close() in defer, when you use this function
func NewGooglePubsub(ctx context.Context, projectID string, opts ...option.ClientOption) *pubSage {
	conn := googlepubsub_adapter.New(ctx, projectID, opts...)
	return newInit(conn)
}

// you can setup the connection by your self
// and use Pubsage as adapter
// call this function and set with your pubsub client
// and Pubsage will take the Publish and Listener
//
// when you use this function, don't call Close()
func GooglePubsubAdapter(c *pubsub.Client) *pubSage {
	conn := googlepubsub_adapter.NewConn(c)
	return newInit(conn)
}
