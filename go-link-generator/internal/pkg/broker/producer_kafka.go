package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// Producer wraps a Sarama SyncProducer and Readiness Checker.
type Producer struct {
	client   sarama.Client
	producer sarama.SyncProducer
	checker  *KafkaReadinessChecker
}

// NewProducer builds a Sarama SyncProducer from the Kafka block in config.
// Return.Successes and Return.Errors are always enabled as required by
// sarama.NewSyncProducerFromClient.
func (k *Kafka) Producer() (*Producer, error) {
	cfg := k.config.Kafka

	scfg := sarama.NewConfig()

	// ── durability ──────────────────────────────────────────
	scfg.Producer.RequiredAcks = requiredAcks(cfg.Producer.RequiredAcks)
	scfg.Producer.Idempotent = true // exactly-once on retry (needs WaitForAll)

	// ── required by SyncProducer ────────────────────────────
	scfg.Producer.Return.Successes = true
	scfg.Producer.Return.Errors = true

	// ── retry / backoff ─────────────────────────────────────
	scfg.Producer.Retry.Max = cfg.Producer.Retry.Max
	scfg.Producer.Retry.Backoff = cfg.Producer.Retry.Backoff

	// ── partitioning ────────────────────────────────────────
	scfg.Producer.Partitioner = partitionerFor(cfg.Producer.Partitioner)

	// ── throughput / message limits ─────────────────────────
	scfg.Producer.Compression = sarama.CompressionSnappy
	scfg.Producer.MaxMessageBytes = cfg.Producer.MaxMessageBytes

	// ── net timeouts ────────────────────────────────────────
	scfg.Net.DialTimeout = cfg.Producer.Timeout.Dial
	scfg.Net.ReadTimeout = cfg.Producer.Timeout.Read
	scfg.Net.WriteTimeout = cfg.Producer.Timeout.Write

	scfg.Net.MaxOpenRequests = 1

	client, err := sarama.NewClient(cfg.Brokers, scfg)
	if err != nil {
		return nil, fmt.Errorf("broker: create kafka client: %w", err)
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("broker: create sync producer: %w", err)
	}

	return &Producer{
		client:   client,
		producer: producer,
		checker:  NewKafkaReadinessChecker(client),
	}, nil
}

// Check validates that the underlying Kafka Producer is healthy.
func (p *Producer) Check(ctx context.Context) error {
	if !p.checker.isHealthy() {
		return fmt.Errorf("kafka producer broker is not healthy or unreachable")
	}
	return nil
}

// KafkaReadinessChecker wraps the Sarama client and caches the health state
// to prevent spamming the Kafka cluster with metadata requests.
type KafkaReadinessChecker struct {
	client sarama.Client

	mu        sync.RWMutex
	isReady   bool
	lastCheck time.Time
	cacheTTL  time.Duration
}

// NewKafkaReadinessChecker creates a new checker with a 10-second cache TTL.
func NewKafkaReadinessChecker(client sarama.Client) *KafkaReadinessChecker {
	return &KafkaReadinessChecker{
		client:   client,
		cacheTTL: 10 * time.Second, // Cache the result for 10 seconds
	}
}

// isHealthy performs the actual check or returns the cached result.
func (c *KafkaReadinessChecker) isHealthy() bool {
	c.mu.RLock()
	// If the cache is still valid, return the cached status
	if time.Since(c.lastCheck) < c.cacheTTL {
		status := c.isReady
		c.mu.RUnlock()
		return status
	}
	c.mu.RUnlock()

	// Cache expired. We need to check Kafka. Lock for writing.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check pattern in case another goroutine updated it while we waited for the write lock
	if time.Since(c.lastCheck) < c.cacheTTL {
		return c.isReady
	}

	// 1. Check if the client is completely closed
	if c.client.Closed() {
		c.updateCache(false)
		return false
	}

	// 2. Force a metadata refresh. If Kafka is unreachable, this fails.
	err := c.client.RefreshMetadata()
	if err != nil {
		log.Printf("Kafka readiness check failed: %v", err)
		c.updateCache(false)
		return false
	}

	// 3. Ensure the client actually recognizes active brokers
	if len(c.client.Brokers()) == 0 {
		log.Println("Kafka readiness check failed: no active brokers found")
		c.updateCache(false)
		return false
	}

	c.updateCache(true)
	return true
}

// updateCache helper to set the state and timestamp
func (c *KafkaReadinessChecker) updateCache(status bool) {
	c.isReady = status
	c.lastCheck = time.Now()
}

type SendConfig struct {
	Header    map[string]string
	Metadata  any
	Offset    int64
	Partition int32
}

type SendOption func(*SendConfig)

func WithHeader(k, v string) SendOption {
	return func(c *SendConfig) {
		if c.Header == nil {
			c.Header = make(map[string]string)
		}
		c.Header[k] = v
	}
}

func WithMetadata(metadata any) SendOption {
	return func(c *SendConfig) {
		c.Metadata = metadata
	}
}

func WithOffset(offset int64) SendOption {
	return func(c *SendConfig) {
		c.Offset = offset
	}
}

func WithPartition(partition int32) SendOption {
	return func(c *SendConfig) {
		c.Partition = partition
	}
}

// BatchMessage represents a single message in a batch operation.
type BatchMessage struct {
	Topic string
	Key   string
	Value any
	Opts  []SendOption
}

// Produce sends a single message to a specific topic, optionally configured with headers/partition/metadata.
// It blocks until the broker acknowledges it.
func (p *Producer) Produce(ctx context.Context, topic string, key string, value any, opts ...SendOption) error {
	m, err := buildMessage(topic, key, value, opts...)
	if err != nil {
		return err
	}

	if _, _, err = p.producer.SendMessage(m); err != nil {
		return fmt.Errorf("broker: producer send: %w", err)
	}

	return nil
}

// ProduceBatch sends multiple messages simultaneously.
func (p *Producer) ProduceBatch(ctx context.Context, messages []BatchMessage) error {
	if len(messages) == 0 {
		return nil
	}

	producerMessages := make([]*sarama.ProducerMessage, 0, len(messages))
	for _, msg := range messages {
		m, err := buildMessage(msg.Topic, msg.Key, msg.Value, msg.Opts...)
		if err != nil {
			return err
		}
		producerMessages = append(producerMessages, m)
	}

	if err := p.producer.SendMessages(producerMessages); err != nil {
		return fmt.Errorf("broker: producer send batch: %w", err)
	}

	return nil
}

// buildMessage encapsulates the logic mapping a generic payload and options into a sarama.ProducerMessage.
func buildMessage(topic string, key string, value any, opts ...SendOption) (*sarama.ProducerMessage, error) {
	var sendConfig SendConfig
	for _, opt := range opts {
		opt(&sendConfig)
	}

	var val []byte
	switch v := value.(type) {
	case []byte:
		val = v
	case string:
		val = []byte(v)
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("broker: marshal message value: %w", err)
		}
		val = encoded
	}

	var headers []sarama.RecordHeader
	if len(sendConfig.Header) > 0 {
		headers = make([]sarama.RecordHeader, 0, len(sendConfig.Header))
		for k, v := range sendConfig.Header {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
	}

	return &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(val),
		Headers:   headers,
		Metadata:  sendConfig.Metadata,
		Offset:    sendConfig.Offset,
		Partition: sendConfig.Partition,
	}, nil
}

// PublishBatch sends multiple messages in a single call.
// All messages succeed or the call returns on the first error.
func (p *Producer) PublishBatch(messages []*sarama.ProducerMessage) error {
	if err := p.producer.SendMessages(messages); err != nil {
		return fmt.Errorf("broker: producer send batch: %w", err)
	}
	return nil
}

// Close flushes any buffered messages and closes the underlying producer.
func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("broker: producer close: %w", err)
	}
	if err := p.client.Close(); err != nil {
		return fmt.Errorf("broker: client close: %w", err)
	}
	return nil
}

func partitionerFor(s string) sarama.PartitionerConstructor {
	switch s {
	case "random":
		return sarama.NewRandomPartitioner
	case "roundrobin":
		return sarama.NewRoundRobinPartitioner
	case "hash":
		return sarama.NewHashPartitioner
	default:
		return sarama.NewHashPartitioner
	}
}

func requiredAcks(s string) sarama.RequiredAcks {
	switch s {
	case "none":
		return sarama.NoResponse
	case "local":
		return sarama.WaitForLocal
	case "all":
		return sarama.WaitForAll
	default:
		return sarama.WaitForAll
	}
}
