package kafka_adapter

import (
	"github.com/Aris-haryanto/pubsage/discourse"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

func parsePublisherToKafka(msg discourse.Publisher) kafka.Message {
	set := kafka.Message{}
	// set general object struct
	set.Topic = msg.Topic
	set.Value = msg.Data
	set.Time = msg.PublishTime
	for k, v := range msg.Attributes {
		set.Headers = append(set.Headers, protocol.Header{
			Key:   k,
			Value: v,
		})
	}

	// ==
	set.Key = msg.Kafka_Key
	set.WriterData = msg.Kafka_WriterData

	return set
}

func parseMessageToPubSage(msg kafka.Message) discourse.Message {
	set := discourse.Message{}
	// set general object struct
	set.Topic = msg.Topic
	set.Data = msg.Value
	set.PublishTime = msg.Time

	attr := map[string][]byte{}
	for _, v := range msg.Headers {
		attr[v.Key] = v.Value
	}
	set.Attributes = attr

	// ==
	set.Kafka_Partition = msg.Partition
	set.Kafka_Offset = msg.Offset
	set.Kafka_HighWaterMark = msg.HighWaterMark
	set.Kafka_Key = msg.Key
	set.Kafka_WriterData = msg.WriterData

	return set
}
