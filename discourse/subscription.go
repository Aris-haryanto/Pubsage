package discourse

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Subscription struct {
	// The topic to read messages from.
	//
	// Available for All
	Topic string

	// when you set true, this will ignoring return true/false from your logic function
	// and will be auto Ack/Commit your message whatever happens
	//
	// in Kafka Only used when GroupID is set
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
