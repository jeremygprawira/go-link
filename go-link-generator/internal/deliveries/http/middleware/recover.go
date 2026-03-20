package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
)

// RecoverMiddleware returns a recovery middleware that logs panics.
func RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Instance.Error(c.Request().Context(), "panic recovered",
						logger.Any("panic", r),
					)
					err = c.JSON(http.StatusInternalServerError, map[string]string{
						"message": "internal server error",
					})
				}
			}()
			return next(c)
		}
	}
}
