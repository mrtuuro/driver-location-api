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

type AddDriversRequest struct {
	Drivers []DriverLocationDTO `json:"drivers" validate:"required,min=1,dive"`
}

type DriverLocationDTO struct {
	DriverId string      `json:"driverId" validate:"required"`
	Location GeoPointDTO `json:"location" validate:"required"`
}

type GeoPointDTO struct {
	Type        string     `json:"type" validate:"required,eq=Point"`
	Coordinates [2]float64 `json:"coordinates" validate:"required,lte=180,gte=-180,dive"`
}

// AddDriversHandler godoc
// @Summary      Bulk-insert driver coordinates
// @Description  Accepts an array of driver GeoJSON points and stores them.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        body  body  handler.AddDriversRequest  true  "Drivers payload"
// @Success      201   {object}  response.SwaggerSuccess
// @Failure      400   {object}  response.SwaggerError
// @Failure      401   {object}  response.SwaggerError
// @Failure      422   {object}  response.SwaggerError
// @Failure      500   {object}  response.SwaggerError
// @Security     InternalAuth
// @Router       /v1/drivers [post]
func AddDriversHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithoutCancel(c.Request().Context())
		var req AddDriversRequest
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

		for _, d := range req.Drivers {
			lat := d.Location.Coordinates[0]
			lon := d.Location.Coordinates[1]
			if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
				return response.RespondError[any](c, apperror.NewAppError(
					code.ErrInvalidLatitudeLongitude,
					errors.New("Drivers locations is not valid."),
					code.GetErrorMessage(code.ErrInvalidLatitudeLongitude),
					))
			}
		}

		drivers := make([]model.DriverLocation, len(req.Drivers))
		for i, d := range req.Drivers {
			drivers[i] = model.DriverLocation{
				DriverID: d.DriverId,
				Location: model.GeoPoint{
					Type:        d.Location.Type,
					Coordinates: d.Location.Coordinates,
				},
			}
		}

		// if len(drivers) == 0 {
		// 	return response.RespondError[any](c, apperror.NewAppError(
		// 		"ERR_NO_DRIVER_PROVIDED",
		// 		errors.New("no drivers provided"),
		// 		"No drivers provided.",
		// 		))
		// }

		if err := app.DriverService.AddDrivers(ctx, drivers); err != nil {
			return response.RespondError[any](c, err)
		}
		return response.RespondSuccess[any](c, code.SuccessDriversCreated, nil)
	}
}
