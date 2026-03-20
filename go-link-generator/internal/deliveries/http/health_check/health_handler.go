package healthcheck

import (
	"net/http"

	"github.com/jeremygprawira/go-link-generator/internal/service"
	"github.com/labstack/echo/v4"
)

// New registers health check routes on the given group.
func New(g *echo.Group, svc *service.Service) {
	g.GET("", handler(svc))
	g.GET("/", handler(svc))
}

func handler(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		status, err := svc.Health.Check(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, status)
	}
}
