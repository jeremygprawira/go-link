package middleware

import (
	"net/http"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/labstack/echo/v4"
)

// APIKeyMiddleware validates the X-API-Key header.
func APIKeyMiddleware(cfg *config.Configuration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-API-Key")
			if key != cfg.Authorization.APIKey {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid or missing API key",
				})
			}
			return next(c)
		}
	}
}
