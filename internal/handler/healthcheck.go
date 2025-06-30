package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mrtuuro/driver-location-api/internal/application"
	"github.com/mrtuuro/driver-location-api/internal/code"
	"github.com/mrtuuro/driver-location-api/internal/response"
)

// HealthcheckHandler godoc
// @Summary      Liveness probe
// @Description  Returns 200 OK with a success envelope; used by load-balancers and orchestrators.
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.SwaggerSuccess
// @Router       /v1/healthz [get]
func HealthcheckHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		return response.RespondSuccess[any](c, code.SuccessHealthCheck, nil)
	}

}
