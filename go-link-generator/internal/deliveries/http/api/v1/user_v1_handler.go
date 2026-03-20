package v1

import (
	"net/http"
	"strings"

	"github.com/jeremygprawira/go-link-generator/internal/config"
	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/response"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/validator"
	"github.com/jeremygprawira/go-link-generator/internal/service"

	"github.com/labstack/echo/v4"
)

type userV1Handler struct {
	service *service.Service
	config  *config.Configuration
}

func NewUserV1(g *echo.Group, svc *service.Service, cfg *config.Configuration) {
	h := &userV1Handler{service: svc, config: cfg}

	users := g.Group("/users")
	users.POST("", h.Create)
	users.GET("/:accountNumber", h.GetByAccountNumber)
}

// Create registers a new user
// @Summary Create New User
// @Description Register a new user with email, phone number, and password. Auto-generates account number.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User Registration Details"
// @Success 201 {object} models.Response{data=models.CreateUserResponse} "User Created Successfully"
// @Failure 400 {object} models.Response "Invalid Input / Validation Error"
// @Failure 409 {object} models.Response "User Already Exists (Email or Phone)"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Router /api/v1/users [post]
func (h *userV1Handler) Create(ctx echo.Context) error {
	var request models.CreateUserRequest
	if err := ctx.Bind(&request); err != nil {
		return response.Error(ctx, err)
	}

	if err := validator.Input(request); err != nil {
		return response.ErrorValidation(ctx, err)
	}

	user, err := h.service.User.Create(ctx.Request().Context(), &request)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, http.StatusCreated, user.CreateUserResponse())
}

// GetUserByAccessToken retrieves user information by access token
// @Summary Get User By Access Token
// @Description Get user information by access token
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.GetUserByAccountNumberResponse} "User Information Retrieved Successfully"
// @Failure 400 {object} models.Response "Invalid Input / Validation Error"
// @Failure 404 {object} models.Response "User Not Found"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Router /api/v1/users/me [get]
// @Security BearerAuth
func (h *userV1Handler) GetByAccountNumber(ctx echo.Context) error {
	user, err := h.service.User.GetByAccountNumber(ctx.Request().Context(), strings.ToLower(ctx.Param("accountNumber")))
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.Success(ctx, http.StatusOK, user.GetUserByAccountNumberResponse())
}
