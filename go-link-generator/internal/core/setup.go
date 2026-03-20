package core

import (
	"context"
	"fmt"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	handler "github.com/jeremygprawira/go-link-generator/internal/deliveries/http"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/database"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/tracer"
	"github.com/jeremygprawira/go-link-generator/internal/repository"
	"github.com/jeremygprawira/go-link-generator/internal/service"
	"github.com/labstack/echo/v4"
)

func Setup(configuration *config.Configuration) (*echo.Echo, error) {
	logger.Initialize(configuration)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

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

	repo := repository.New(db)

	svc := service.New(service.Dependencies{
		Repository: *repo,
		Config:     configuration,
	})
	handler.New(e, svc, configuration)
	return e, nil
}
