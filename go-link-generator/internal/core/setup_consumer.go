package core

import (
	"context"
	"fmt"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	delivery "github.com/jeremygprawira/go-link-generator/internal/deliveries/kafka"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/broker"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/database"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/tracer"
	"github.com/jeremygprawira/go-link-generator/internal/repository"
	"github.com/jeremygprawira/go-link-generator/internal/service"
)

// SetupConsumer wires all core infrastructure and dependency injections specifically
// designed for background message workers like the Kafka consumer.
func SetupConsumer(configuration *config.Configuration) (*broker.Consumer, error) {
	logger.Initialize(configuration)

	tp, err := tracer.Initialize(configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}
	_ = tp

	db, err := database.Connect(configuration)
	if err != nil {
		logger.Instance.Error(context.Background(), "failed to connect to database", logger.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	producer, err := broker.New(configuration).Producer()
	if err != nil {
		logger.Instance.Error(context.Background(), "failed to connect to kafka producer", logger.Error(err))
		return nil, fmt.Errorf("failed to connect to kafka producer: %w", err)
	}

	repo := repository.New(db)

	svc := service.New(service.Dependencies{
		Repository: *repo,
		Config:     configuration,
		Producer:   producer,
	})

	// Multiplexer for incoming streams
	handler := delivery.New(svc, configuration)

	consumer, err := broker.New(configuration).NewConsumer(handler.HandleMessage)
	if err != nil {
		logger.Instance.Error(context.Background(), "failed to construct kafka consumer", logger.Error(err))
		return nil, fmt.Errorf("failed to construct kafka consumer: %w", err)
	}

	return consumer, nil
}
