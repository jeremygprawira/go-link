package v1

import (
	"net/http"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/response"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/validator"
	"github.com/jeremygprawira/go-link-generator/internal/service"

	"github.com/labstack/echo/v4"
)

type urlV1Handler struct {
	service *service.Service
	config  *config.Configuration
}

func NewUrlV1(g *echo.Group, service *service.Service, config *config.Configuration) {
	h := &urlV1Handler{service: service, config: config}

	urls := g.Group("/url")
	urls.POST("", h.Create)
}

// Create generates a new short URL
// @Summary Create New Short URL
// @Description Creates a new shortened URL, optionally with a customized code. Auto-generates a highly-secure Snowflake ID if no code is provided.
// @Tags URLs
// @Accept json
// @Produce json
// @Param request body models.CreateUrlRequest true "URL Creation Details"
// @Success 201 {object} models.Response{data=models.CreateUrlResponse} "URL Created Successfully"
// @Failure 400 {object} models.Response "Invalid Input / Validation Error"
// @Failure 409 {object} models.Response "URL Code Already Exists"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Router /api/v1/url [post]
func (h *urlV1Handler) Create(ctx echo.Context) error {
	var request models.CreateUrlRequest
	if err := ctx.Bind(&request); err != nil {
		return response.Error(ctx, err)
	}

	if err := validator.Input(request); err != nil {
		return response.ErrorValidation(ctx, err)
	}

	url, err := h.service.Url.Create(ctx.Request().Context(), &request)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, http.StatusCreated, url.Response())
}
