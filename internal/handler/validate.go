package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/application"
	"github.com/mrtuuro/driver-location-api/internal/code"
)

func BindAndValidate[T any](app *application.Application, c echo.Context, dest any) (any, error) {
	if err := c.Bind(&dest); err != nil {
		return nil, apperror.NewAppError(code.ErrSystemInternal,
			err,
			code.GetErrorMessage(code.ErrSystemInternal))
	}

	if err := app.E.Validator.Validate(dest); err != nil {
		return nil, apperror.NewAppError(code.ErrValidationFailed,
			err,
			code.GetErrorMessage(code.ErrValidationFailed))
	}

	return dest, nil
}
