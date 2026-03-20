package v1

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/service"

	"github.com/labstack/echo/v4"
)
func New(api *echo.Group, svc *service.Service, cfg *config.Configuration) {
	v1 := api.Group("/v1")
	NewUserV1(v1, svc, cfg)
}
