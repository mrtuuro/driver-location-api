package router

import (
	"fmt"

	"github.com/mrtuuro/driver-location-api/internal/application"
	"github.com/mrtuuro/driver-location-api/internal/handler"
	"github.com/mrtuuro/driver-location-api/internal/middleware"
)

func Register(app *application.Application) {
	app.E.GET("/v1/healthz", handler.HealthcheckHandler(app)).Name = "Healthcheck"

	protected := app.E.Group("/v1")
	protected.Use(middleware.CustomMiddleware(app))

	protected.POST("/drivers", handler.AddDriversHandler(app)).Name = "Add Drivers"
	protected.POST("/drivers/search", handler.SearchDriverHandler(app)).Name = "Search Driver"

}

func PrintRoutes(app *application.Application) {
	fmt.Println("=== ROUTES ===")
	routes := app.E.Routes()
	for _, r := range routes {
		fmt.Printf("%s - [%s]%s\n", r.Name, r.Method, r.Path)
	}
}
