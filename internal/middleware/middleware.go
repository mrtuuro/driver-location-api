package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/application"
	"github.com/mrtuuro/driver-location-api/internal/code"
	"github.com/mrtuuro/driver-location-api/internal/response"
)

func CustomMiddleware(app *application.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			headers := c.Request().Header
			accTokenStr, err := extractBearerToken(headers)
			if err != nil {
				return response.RespondError[any](c, err)
			}

			claims, err := app.TokenManager.ValidateJWT(accTokenStr)
			if err != nil {
				fmt.Println("here 2")
				return response.RespondError[any](c, err)
			}

			if !claims.Authenticated {
				fmt.Println("here")
				return response.RespondError[any](c, apperror.NewAppError(
					code.ErrAuthInvalidCredentials,
					errors.New("Request is not authenticated"),
					code.GetErrorMessage(code.ErrAuthInvalidCredentials),
					))
			}

			return next(c)
		}
	}
}

func extractBearerToken(headers map[string][]string) (string, error) {
	fullAuthToken, ok := headers["Authorization"]
	if !ok {
		return "", apperror.NewAppError(
			code.ErrAuthMissingToken,
			errors.New(code.GetErrorMessage(code.ErrAuthMissingToken)),
			code.GetErrorMessage(code.ErrAuthMissingToken),
			)
	}

	if !strings.HasPrefix(fullAuthToken[0], "Bearer ") {
		return "", apperror.NewAppError(
			code.ErrAuthInvalidProtocol,
			errors.New("Authantication method is not allowed."),
			"Authantication method is wrong or token is missing.",
			)
	}
	accTokenStr := strings.TrimPrefix(fullAuthToken[0], "Bearer ")

	return accTokenStr, nil
}
