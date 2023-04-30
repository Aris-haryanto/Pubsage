package kafka_adapter

import "github.com/segmentio/kafka-go"

func New(w *kafka.Writer) *KafkaAdapter {
	return &KafkaAdapter{writer: w}
}
