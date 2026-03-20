package response

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/errorc"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/logger"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/stringc"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/labstack/echo/v4"
)
func Error(ctx echo.Context, err error, message ...string) error {
	var response models.Response

	if err == nil {
		errorcResponse := errorc.GetResponse(errorc.ErrorInternalServer)
		response = models.Response{
			Code:    errorcResponse.Code,
			Status:  errorcResponse.Status,
			Message: errorcResponse.Message,
		}
	} else {
		errorcResponse := errorc.GetResponse(err)

		var finalMessage string
		switch len(message) {
		case 0:
			finalMessage = errorcResponse.Message
		case 1:
			finalMessage = message[0]
		default:
			finalMessage = fmt.Sprintf(message[0], stringc.SlicesToInterfaces(message[1:])...)
		}

		response = models.Response{
			Code:    errorcResponse.Code,
			Status:  errorcResponse.Status,
			Message: finalMessage,
		}

		// Log the error securely
		var httpErr *errorc.HTTPError
		var fwErr *echo.HTTPError

		if errors.As(err, &httpErr) {
			// It's a structured AppError
			errType := "AppError"
			var errDetails any

			if httpErr.Internal() != nil {
				errDetails = httpErr.Internal().Error()
			}

			if errorcResponse.Code == http.StatusInternalServerError {
				errType = "InternalError"
			}
			logger.AddError(ctx.Request().Context(), &logger.ErrorContext{
				Type:    errType,
				Message: httpErr.Error(), // Always log the public/predefined message
				Details: errDetails,      // Log the actual internal error here if it exists
			})
		} else if errors.As(err, &fwErr) {
			// Echo framework error (e.g., binding failure)
			var errDetails any
			if fwErr.Internal != nil {
				errDetails = fwErr.Internal.Error()
			}

			logger.AddError(ctx.Request().Context(), &logger.ErrorContext{
				Type:    "FrameworkError",
				Message: fmt.Sprintf("%v", fwErr.Message),
				Details: errDetails,
			})
		} else {
			// It's a raw, unknown error
			logger.AddError(ctx.Request().Context(), &logger.ErrorContext{
				Type:    "RawError",
				Message: err.Error(),
			})
		}
	}
	requestID, _ := ctx.Get("X-Request-ID").(string)
	timestamp, _ := ctx.Get("X-Timestamp").(string)

	if requestID == "" {
		requestID = uuid.New().String()
	}
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	response.Metadata = models.Metadata{
		RequestId: requestID,
		Timestamp: timestamp,
	}
	return ctx.JSON(response.Code, response)
}
func ErrorValidation(ctx echo.Context, errors interface{}) error {
	requestID, _ := ctx.Get("X-Request-ID").(string)
	timestamp, _ := ctx.Get("X-Timestamp").(string)

	if requestID == "" {
		requestID = uuid.New().String()
	}
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	errorcResponse := errorc.GetResponse(errorc.ErrorValidation)
	response := models.Response{
		Code:    errorcResponse.Code,
		Status:  errorcResponse.Status,
		Message: errorcResponse.Message,
		Metadata: models.Metadata{
			RequestId: requestID,
			Timestamp: timestamp,
		},
	}

	if data, ok := errors.(*multierror.Error); ok {
		response.Errors = data.Errors
	}
	return ctx.JSON(response.Code, response)
}
