package middleware

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// Middleware holds middleware setup.
type Middleware struct {
	e   *echo.Echo
	cfg *config.Configuration
}

// New creates a new Middleware instance.
func New(e *echo.Echo, cfg *config.Configuration) *Middleware {
	return &Middleware{e: e, cfg: cfg}
}

// Default registers essential middleware (recover, request ID, CORS, logger).
func (m *Middleware) Default(cfg *config.Configuration) {
	m.e.Use(echomiddleware.Recover())
	m.e.Use(echomiddleware.RequestID())
	m.e.Use(CORSMiddleware(cfg))
	m.e.Use(m.LoggingMiddleware(logger.Instance))
}
