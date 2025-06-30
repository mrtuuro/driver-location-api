// @title          Driver Location API
// @version        1.0
// @description    Stores driver coordinates and returns nearest drivers
// @BasePath       /v1

// @securityDefinitions.apikey InternalAuth
// @in header
// @name Authorization
// @description  Internal calls only.  Format: "Bearer <token>"
package main

import (
	"log"
	"os"

	"github.com/mrtuuro/driver-location-api/internal/application"
	"github.com/mrtuuro/driver-location-api/internal/config"
	"github.com/mrtuuro/driver-location-api/internal/db"
	"github.com/mrtuuro/driver-location-api/internal/repository"
	"github.com/mrtuuro/driver-location-api/internal/router"
	"github.com/mrtuuro/driver-location-api/internal/service"
)

func main() {
	cfg := config.NewConfig()

	mongoClient, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("err connecting db: %v", err)
		os.Exit(1)
	}

	var (
		// COLLECTION VARIABLES
		driverColl = db.GetCollection(mongoClient, cfg.DatabaseName, cfg.CollectionName)

		// REPOSITORY INITIALIZERS
		driverRepo = repository.NewMongoDriverRepository(driverColl)

		// SERVICE INIT
		driverSvc = service.NewDriverService(driverRepo)
	)

	app := application.NewApp(cfg, driverSvc)

	router.Register(app)
	router.PrintRoutes(app)

	app.Run(cfg.Port)

}
