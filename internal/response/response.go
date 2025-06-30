package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/code"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse[T any] struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Code    string    `json:"code,omitempty"`
	Data    *T        `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

func RespondSuccess[T any](c echo.Context, statusCode string, data *T) error {
	return c.JSON(code.GetStatusCode(statusCode), APIResponse[T]{
		Success: true,
		Code:    statusCode,
		Message: code.GetSuccessMessage(statusCode),
		Data:    data,
	})
}

func RespondError[T any](c echo.Context, err error) error {
	appErr, ok := err.(*apperror.AppError)
	if !ok {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	}

	status := code.GetStatusCode(appErr.Code)

	fmt.Printf("\n%v", appErr)
	return c.JSON(status, APIResponse[T]{
		Success: false,
		Error: &APIError{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
