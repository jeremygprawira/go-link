package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// MessageHandler is the function your application provides to process
// a single message. Return a non-nil error to signal a processing failure.
type MessageHandler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// Consumer wraps a Sarama ConsumerGroup and owns its run-loop.
type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler *consumerGroupHandler
	cancel  context.CancelFunc
	done    chan struct{}
}

// NewConsumer builds and returns a Consumer ready to call Start() on.
func (k *Kafka) NewConsumer(fn MessageHandler) (*Consumer, error) {
	cfg := k.config.Kafka

	scfg := sarama.NewConfig()
	scfg.Version = sarama.V3_4_0_0

	// --- offset strategy ---
	switch cfg.Consumer.InitialOffset {
	case "oldest":
		scfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	default:
		scfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	}
	scfg.Consumer.Offsets.AutoCommit.Enable = false // we commit manually after processing

	// --- group membership ---
	scfg.Consumer.Group.Session.Timeout = cfg.Consumer.Timeout.Session
	scfg.Consumer.Group.Heartbeat.Interval = cfg.Consumer.HeartbeatInterval
	scfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	scfg.Consumer.Group.Rebalance.Timeout = cfg.Consumer.Timeout.Rebalance

	// --- retries & timeouts ---
	scfg.Consumer.Retry.Backoff = cfg.Consumer.Retry.Backoff
	scfg.Net.DialTimeout = cfg.Consumer.Timeout.Dial
	scfg.Net.ReadTimeout = cfg.Consumer.Timeout.Read
	scfg.Net.WriteTimeout = cfg.Consumer.Timeout.Write

	// --- fetch tuning ---
	scfg.Consumer.Fetch.Min = 1
	scfg.Consumer.Fetch.Default = 1 << 20 // 1 MiB
	scfg.Consumer.MaxWaitTime = 500 * time.Millisecond

	group, err := sarama.NewConsumerGroup(k.config.Kafka.Brokers, cfg.Consumer.GroupID, scfg)
	if err != nil {
		return nil, fmt.Errorf("kafka: create consumer group: %w", err)
	}

	return &Consumer{
		group:   group,
		topics:  []string{cfg.Topics.Link, cfg.Topics.LinkDLQ},
		handler: &consumerGroupHandler{fn: fn},
		done:    make(chan struct{}),
	}, nil
}

// Start launches the consume loop in a background goroutine.
// It returns immediately; call Stop() to shut down cleanly.
func (c *Consumer) Start(ctx context.Context) {
	ctx, c.cancel = context.WithCancel(ctx)

	go func() {
		defer close(c.done)
		for {
			// Consume re-joins the group after every rebalance automatically.
			if err := c.group.Consume(ctx, c.topics, c.handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				// Surface transient errors without dying — the loop retries.
				// In production, feed this into your metrics / alert pipeline.
				log.Printf("kafka: consumer error: %v", err)
			}
			if ctx.Err() != nil {
				return // context cancelled — clean exit
			}
		}
	}()
}

// Stop signals the run-loop to exit and waits for it to finish.
func (c *Consumer) Stop() error {
	c.cancel()
	<-c.done
	return c.group.Close()
}

// -----------------------------------------------------------------
// consumerGroupHandler implements sarama.ConsumerGroupHandler.
// -----------------------------------------------------------------

type consumerGroupHandler struct {
	fn    MessageHandler
	ready chan struct{}
}

// Setup is called by Sarama at the start of each new session (after every
// rebalance). Signal readiness so callers can wait before producing.
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	if h.ready == nil {
		h.ready = make(chan struct{})
	}
	close(h.ready)
	return nil
}

// Cleanup is called at the end of a session. Release any per-session state.
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.ready = make(chan struct{}) // reset for next session
	return nil
}

// ConsumeClaim is the hot path — one goroutine per topic/partition claim.
// We commit only after the handler acknowledges the message successfully.
func (h *consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil // partition revoked
			}
			if err := h.fn(session.Context(), msg); err != nil {
				// Log and continue — dead-letter queue logic belongs in fn itself.
				log.Printf("kafka: handler error topic=%s partition=%d offset=%d: %v",
					msg.Topic, msg.Partition, msg.Offset, err)
				continue
			}
			// Mark only after a successful handler return.
			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// WaitReady blocks until the consumer has completed its first rebalance
// and is actively consuming. Useful in tests and integration harnesses.
func (c *Consumer) WaitReady(ctx context.Context) error {
	select {
	case <-c.handler.ready:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
