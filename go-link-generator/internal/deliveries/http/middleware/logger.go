package middleware

import (
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// LoggerMiddleware returns an Echo logger middleware using the app logger.
func LoggerMiddleware() echo.MiddlewareFunc {
	return echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogMethod: true,
		LogStatus: true,
		LogError:  true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			if v.Error != nil {
				logger.Instance.Error(c.Request().Context(), "request",
					logger.String("method", v.Method),
					logger.String("uri", v.URI),
					logger.Int("status", v.Status),
					logger.Error(v.Error),
				)
			} else {
				logger.Instance.Info(c.Request().Context(), "request",
					logger.String("method", v.Method),
					logger.String("uri", v.URI),
					logger.Int("status", v.Status),
				)
			}
			return nil
		},
	})
}
