package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/core"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/graceful"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
)

func main() {
	// Run the application and handle any startup errors
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application failed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Initialize configuration first (before any dependencies)
	ctx := context.Background()

	cfg, err := config.Initialize(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to initialize configuration: %w", err))
	}

	// 2. Setup application core, from logger, to echo server, database, etc.
	server, err := core.Setup(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup application: %w", err)
	}

	port := fmt.Sprintf(":%d", cfg.Application.Port)

	// 3. Define processes to manage
	processes := map[string]graceful.Process{
		"http-server": graceful.NewEchoProcess(server, port),
		"cleanup":     graceful.NewFuncProcess(core.Teardown),
	}

	// 4. Configure graceful shutdown
	shutdownTimeout := 10 * time.Second
	logAdapter := graceful.NewLoggerAdapter(logger.Instance, ctx)

	// 5. Run with graceful lifecycle management
	graceful.Graceful(processes,
		graceful.WithTimeout(shutdownTimeout),
		graceful.WithLogger(logAdapter),
		graceful.WithStartupHook(func() {
			logger.Instance.Info(ctx, "All processes started successfully",
				logger.String("env", cfg.Application.Environment),
				logger.String("addr", port),
			)
		}),
		graceful.WithShutdownHook(func() {
			logger.Instance.Info(ctx, "Beginning graceful shutdown",
				logger.String("timeout", shutdownTimeout.String()),
			)
		}),
	)

	logger.Instance.Info(ctx, "Server shutdown completed successfully")
	return nil
}
