package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
	"github.com/jeremygprawira/go-link-generator/internal/service"
)

// Handler delegates incoming Kafka messages to the appropriate business logic service.
type Handler struct {
	svc *service.Service
	cfg *config.Configuration
}

// New constructs a new kafka routes multiplexer.
func New(svc *service.Service, cfg *config.Configuration) *Handler {
	return &Handler{
		svc: svc,
		cfg: cfg,
	}
}

// HandleMessage is invoked concurrently for each incoming message belonging to the consumer group claims.
func (h *Handler) HandleMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	logger.Instance.Info(ctx, "processing kafka message",
		logger.String("topic", msg.Topic),
		logger.String("partition", fmt.Sprintf("%d", msg.Partition)),
	)

	// Route based on topics
	switch msg.Topic {
	case h.cfg.Kafka.Topics.Link:
		// process logic here using h.svc (e.g. h.svc.Url)
	case h.cfg.Kafka.Topics.LinkDLQ:
		// DLQ logic mapped
	default:
		logger.Instance.Warn(ctx, "unrecognized topic received", logger.String("topic", msg.Topic))
	}

	return nil
}
