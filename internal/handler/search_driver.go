package handler

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/driver-location-api/internal/apperror"
	"github.com/mrtuuro/driver-location-api/internal/application"
	"github.com/mrtuuro/driver-location-api/internal/code"
	"github.com/mrtuuro/driver-location-api/internal/model"
	"github.com/mrtuuro/driver-location-api/internal/response"
)

type SearchDriverRequest struct {
	Location GeoPointDTO `json:"location" validate:"required"`
	Radius   float64     `json:"radius" validate:"required,gt=0"`
	Limit    int         `json:"limit"`
}

// SearchDriverHandler godoc
// @Summary      Find nearest drivers
// @Description  Returns drivers ordered by distance; distance (metres) is pre-calculated.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        body  body  handler.SearchDriverRequest  true  "Search parameters"
// @Success      200   {object}  response.SwaggerSuccess        "List of DriverWithDistance"
// @Failure      400   {object}  response.SwaggerError
// @Failure      401   {object}  response.SwaggerError
// @Failure      404   {object}  response.SwaggerError
// @Failure      500   {object}  response.SwaggerError
// @Security     InternalAuth
// @Router       /v1/drivers/search [post]
func SearchDriverHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithoutCancel(c.Request().Context())
		var req SearchDriverRequest
		if err := c.Bind(&req); err != nil {
			return response.RespondError[any](c, apperror.NewAppError(
				code.ErrInvalidJSON,
				err,
				code.GetErrorMessage(code.ErrInvalidJSON),
				))
		}

		if err := c.Validate(&req); err != nil {
			return response.RespondError[any](c, apperror.NewAppError(
				code.ErrValidationFailed,
				err,
				code.GetErrorMessage(code.ErrValidationFailed),
				))
		}

		lat := req.Location.Coordinates[0]
		lon := req.Location.Coordinates[1]
		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			return response.RespondError[any](c, apperror.NewAppError(
				code.ErrInvalidLatitudeLongitude,
				errors.New("Drivers locations is not valid."),
				code.GetErrorMessage(code.ErrInvalidLatitudeLongitude),
				))
		}

		userPoint := &model.GeoPoint{
			Type:        req.Location.Type,
			Coordinates: req.Location.Coordinates,
		}

		driverLocs, err := app.DriverService.SearchDriver(ctx, userPoint, req.Radius, req.Limit)
		if err != nil {
			return response.RespondError[any](c, err)
		}

		return response.RespondSuccess[[]model.DriverWithDistance](c, code.SuccessOperationCompleted, &driverLocs)
	}
}
