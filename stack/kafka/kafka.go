package kafka_adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/Aris-haryanto/pubsage/discourse"
	"github.com/segmentio/kafka-go"
)

type KafkaAdapter struct {
	writer *kafka.Writer
}

func (k *KafkaAdapter) Close() error {
	if err := k.writer.Close(); err != nil {
		return fmt.Errorf("[Pubsage-Kafka] error when close kafka writer connection - %v", err)
	}

	return nil
}

func (k *KafkaAdapter) Publish(ctx context.Context, message discourse.Publisher) error {

	msg := parsePublisherToKafka(message)

	err := k.writer.WriteMessages(ctx, msg)

	return err
}

func (k *KafkaAdapter) Listener(ctx context.Context, fn func(discourse.Message) error, cfgSubcription discourse.Subscription) error {
	reader := kafka.NewReader(mapReaderConfig(cfgSubcription))
	defer reader.Close()

	var err error
	if cfgSubcription.AutoCommit {
		err = fetchWithAutoCommit(ctx, reader, fn)
	} else {
		err = fetchWithManualCommit(ctx, reader, fn)
	}

	return err

}

func fetchWithAutoCommit(ctx context.Context, reader *kafka.Reader, fn func(discourse.Message) error) error {
	var (
		getReader   kafka.Message
		readerError error
	)
	for {
		getReader, readerError = reader.ReadMessage(ctx)
		if readerError != nil {
			break
		}

		message := parseMessageToPubSage(getReader)

		if errFn := fn(message); errFn != nil {
			fmt.Printf("[Pubsage-Kafka] message with key %s - Error: %s\n", message.Kafka_Key, errFn)
			continue
		}
	}

	return readerError
}

func fetchWithManualCommit(ctx context.Context, reader *kafka.Reader, fn func(discourse.Message) error) error {
	var (
		getReader   kafka.Message
		readerError error
	)

	for {
		getReader, readerError = reader.FetchMessage(ctx)
		if readerError != nil {
			break
		}

		message := parseMessageToPubSage(getReader)

		if errFn := fn(message); errFn != nil {
			fmt.Printf("[Pubsage-Kafka] skip commiting message with key %s - Error: %s\n", message.Kafka_Key, errFn)
			continue
		}

		// commit if no error
		reader.CommitMessages(ctx, getReader)
	}

	return readerError
}

func mapReaderConfig(cfgReader discourse.Subscription) kafka.ReaderConfig {
	var commitInterval time.Duration
	if cfgReader.AutoCommit {
		commitInterval = time.Second
	}

	return kafka.ReaderConfig{
		Topic:                  cfgReader.Topic,
		Brokers:                cfgReader.Kafka_Brokers,
		GroupID:                cfgReader.Kafka_GroupID,
		GroupTopics:            cfgReader.Kafka_GroupTopics,
		Partition:              cfgReader.Kafka_Partition,
		Dialer:                 cfgReader.Kafka_Dialer,
		QueueCapacity:          cfgReader.Kafka_QueueCapacity,
		MinBytes:               cfgReader.Kafka_MinBytes,
		MaxBytes:               cfgReader.Kafka_MaxBytes,
		MaxWait:                cfgReader.Kafka_MaxWait,
		ReadBatchTimeout:       cfgReader.Kafka_ReadBatchTimeout,
		ReadLagInterval:        cfgReader.Kafka_ReadLagInterval,
		GroupBalancers:         cfgReader.Kafka_GroupBalancers,
		HeartbeatInterval:      cfgReader.Kafka_HeartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: cfgReader.Kafka_PartitionWatchInterval,
		WatchPartitionChanges:  cfgReader.Kafka_WatchPartitionChanges,
		SessionTimeout:         cfgReader.Kafka_SessionTimeout,
		RebalanceTimeout:       cfgReader.Kafka_RebalanceTimeout,
		JoinGroupBackoff:       cfgReader.Kafka_JoinGroupBackoff,
		RetentionTime:          cfgReader.Kafka_RetentionTime,
		StartOffset:            cfgReader.Kafka_StartOffset,
		ReadBackoffMin:         cfgReader.Kafka_ReadBackoffMin,
		ReadBackoffMax:         cfgReader.Kafka_ReadBackoffMax,
		Logger:                 cfgReader.Kafka_Logger,
		ErrorLogger:            cfgReader.Kafka_ErrorLogger,
		IsolationLevel:         cfgReader.Kafka_IsolationLevel,
		MaxAttempts:            cfgReader.Kafka_MaxAttempts,
		OffsetOutOfRangeError:  cfgReader.Kafka_OffsetOutOfRangeError,
	}
}
