package http

import (
	"net/http"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	v1 "github.com/jeremygprawira/go-link-generator/internal/deliveries/http/api/v1"
	healthcheck "github.com/jeremygprawira/go-link-generator/internal/deliveries/http/health_check"
	"github.com/jeremygprawira/go-link-generator/internal/deliveries/http/middleware"
	"github.com/jeremygprawira/go-link-generator/internal/service"

	"github.com/labstack/echo/v4"
	_ "github.com/jeremygprawira/go-link-generator/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)
// @title Go-Link-Generator API
// @version 1.0
// @description Auto-generated API documentation for go-link-generator.
// @host localhost:8080
// @BasePath /api
// @schemes http https
func New(e *echo.Echo, svc *service.Service, cfg *config.Configuration) {
	mw := middleware.New(e, cfg)
	mw.Default(cfg)

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "go-link-generator is running 🚀")
	})
	e.GET("/docs", func(ctx echo.Context) error {
		return ctx.File("docs/api-docs.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Health
	health := e.Group("/health")
	healthcheck.New(health, svc)

	// API
	api := e.Group("/api")
	api.Use(middleware.APIKeyMiddleware(cfg))

	// V1
	v1.New(api, svc, cfg)
}
