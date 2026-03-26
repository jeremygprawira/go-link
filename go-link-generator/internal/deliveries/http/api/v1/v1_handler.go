package v1

import (
	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/service"

	"github.com/labstack/echo/v4"
)

func New(api *echo.Group, service *service.Service, config *config.Configuration) {
	v1 := api.Group("/v1")
	NewUrlV1(v1, service, config)
}
