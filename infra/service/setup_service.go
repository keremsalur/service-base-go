package service

import (
	database "service-base-go/infra/db"
	"service-base-go/infra/route"
	"service-base-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ServiceConfig struct {
	ServiceName string
	Routes      route.Routes
}

func GetServicesConfigs() []ServiceConfig {
	return []ServiceConfig{
		{
			ServiceName: "Health Service",
			Routes:      route.NewHealthCheckRoutes(),
		},

		{
			ServiceName: "Project Service",
			Routes:      route.NewProjectRoutes(),
		},
		/*
			{
				ServiceName: "Company Service",
				Routes:      routes.NewCompanyRoutes(),
			},
			{
				ServiceName: "Location Service",
				Routes:      routes.NewLocationRoutes(),
			},
			{
				ServiceName: "Product Type Service",
				Routes:      routes.NewProductTypeRoutes(),
			},
			{
				ServiceName: "Type Service",
				Routes:      routes.NewTypeRoutes(),
			},
			{
				ServiceName: "Plate Service",
				Routes:      routes.NewPlateRoutes(),
			},
			{
				ServiceName: "Cut Service",
				Routes:      routes.NewCutRoutes(),
			},
		*/
	}
}

func SetupServices(app *fiber.App, db database.Database) error {
	database := db.GetDB()
	logData := logger.GetGlobalLogData()
	logData["class"] = "SetupService"
	// Tüm servisleri dinamik olarak başlat
	for _, config := range GetServicesConfigs() {
		err := setupService(app, config, database)
		if err != nil {
			logData["error"] = err
			logger.PushLog(logger.GetLogger(), zerolog.FatalLevel, config.ServiceName+" servisi başlatılamadı", logData)
			return err
		}
		logger.PushLog(logger.GetLogger(), zerolog.InfoLevel, config.ServiceName+" servisi başlatıldı.", logData)
	}
	return nil
}

func setupService(app *fiber.App, config ServiceConfig, db *gorm.DB) error {
	config.Routes.RegisterRoutes(app, db)
	return nil
}
