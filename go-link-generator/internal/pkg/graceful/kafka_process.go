package graceful

import (
	"context"

	"github.com/jeremygprawira/go-link-generator/internal/pkg/broker"
)

// KafkaProcess wraps a broker.Consumer into the graceful lifecycle Process.
type KafkaProcess struct {
	consumer *broker.Consumer
}

// NewKafkaProcess constructs the lifecycle process given the configured consumer.
func NewKafkaProcess(consumer *broker.Consumer) *KafkaProcess {
	return &KafkaProcess{consumer: consumer}
}

// Start will launch the consumer loop and block until it's ready, and then
// block forever until ctx is canceled (meaning application is shutting down).
func (k *KafkaProcess) Start(ctx context.Context) error {
	k.consumer.Start(ctx)

	// Wait until consumer correctly establishes the session and is ready to ingest
	if err := k.consumer.WaitReady(ctx); err != nil {
		return err
	}

	// Block here so graceful manager considers process "running"
	<-ctx.Done()
	return nil
}

// Stop executes the blocking shutdown procedure of the underlying consumer group.
func (k *KafkaProcess) Stop(_ context.Context) error {
	return k.consumer.Stop()
}
