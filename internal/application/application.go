package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/mrtuuro/driver-location-api/docs"
	"github.com/mrtuuro/driver-location-api/internal/config"
	"github.com/mrtuuro/driver-location-api/internal/service"
	"github.com/mrtuuro/driver-location-api/internal/token"
	"github.com/mrtuuro/driver-location-api/internal/validator"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Application struct {
	Cfg *config.Config

	E            *echo.Echo
	TokenManager *token.TokenManager

	// SERVICES
	DriverService service.DriverService
}

func NewApp(cfg *config.Config, driverService service.DriverService) *Application {
	app := &Application{}
	tm := token.NewTokenManager(cfg.SecretKey)

	app.Cfg = cfg
	app.DriverService = driverService
	app.E = setupEcho()
	app.TokenManager = tm

	return app
}

func setupEcho() *echo.Echo {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Validator = validator.NewCustomValidator()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.ERROR)
	return e
}

func (app *Application) Run(port string) {
	ctx, stop := signal.NotifyContext(app.Cfg.Ctx, os.Interrupt)
	defer stop()

	go func() {
		if err := app.E.Start(port); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			app.E.Logger.Fatal("Shutting down the server!")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(app.Cfg.Ctx, 10*time.Second)
	defer cancel()
	if err := app.E.Shutdown(ctx); err != nil {
		app.E.Logger.Fatal(err)
	}
}
