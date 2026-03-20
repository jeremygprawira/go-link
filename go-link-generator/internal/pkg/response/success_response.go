package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jeremygprawira/go-link-generator/internal/models"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/numberc"
	"github.com/jeremygprawira/go-link-generator/internal/pkg/stringc"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)
func Success(ctx echo.Context, code int, data interface{}, message ...string) error {
	requestID, _ := ctx.Get("X-Request-ID").(string)
	timestamp, _ := ctx.Get("X-Timestamp").(string)

	if requestID == "" {
		requestID = uuid.New().String()
	}
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	var finalMessage string
	switch len(message) {
	case 0:
		finalMessage = "Request has been successfully processed."
	case 1:
		finalMessage = message[0]
	default:
		finalMessage = fmt.Sprintf(message[0], stringc.SlicesToInterfaces(message[1:])...)
	}

	resp := models.Response{
		Code:    code,
		Status:  "OK",
		Message: finalMessage,
		Data:    data,
		Metadata: models.Metadata{
			RequestId: requestID,
			Timestamp: timestamp,
		},
	}
	return ctx.JSON(code, resp)
}
func SuccessList(ctx echo.Context, code int, message string, data interface{}) error {
	requestID, _ := ctx.Get("X-Request-ID").(string)
	timestamp, _ := ctx.Get("X-Timestamp").(string)

	if requestID == "" {
		requestID = uuid.New().String()
	}
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	if message == "" {
		message = "Request has been successfully processed."
	}

	length, err := numberc.LengthOf(data)
	if err != nil {
		errResp := models.Response{
			Code:    http.StatusInternalServerError,
			Status:  "ERROR",
			Message: "Failed to get length of data.",
			Metadata: models.Metadata{
				RequestId: requestID,
				Timestamp: timestamp,
			},
		}
		return ctx.JSON(http.StatusInternalServerError, errResp)
	}

	resp := models.Response{
		Code:    code,
		Status:  "OK",
		Message: message,
		Data:    data,
		Metadata: models.Metadata{
			RequestId: requestID,
			Timestamp: timestamp,
			TotalRows: length,
		},
	}
	return ctx.JSON(http.StatusOK, resp)
}
func SuccessPagination(ctx echo.Context, code int, message string, pagination models.PaginationOutput, data interface{}) error {
	requestID, _ := ctx.Get("X-Request-ID").(string)
	timestamp, _ := ctx.Get("X-Timestamp").(string)

	if requestID == "" {
		requestID = uuid.New().String()
	}
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	if message == "" {
		message = "Request has been successfully processed."
	}

	resp := models.Response{
		Code:       code,
		Status:     "OK",
		Message:    message,
		Data:       data,
		Pagination: &pagination,
		Metadata: models.Metadata{
			RequestId: requestID,
			Timestamp: timestamp,
		},
	}
	return ctx.JSON(http.StatusOK, resp)
}
