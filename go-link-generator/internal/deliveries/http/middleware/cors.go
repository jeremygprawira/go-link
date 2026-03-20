package middleware

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSMiddleware configures CORS for the Echo instance.
func CORSMiddleware(cfg *config.Configuration) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: cfg.CORS.HeadersAllowed,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	})
}
