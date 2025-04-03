package route

import (
	"service-base-go/app/healthcheck"
	"service-base-go/pkg/handler"
	"service-base-go/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthcheckRoutes struct {
}

func NewHealthCheckRoutes() Routes {
	return &HealthcheckRoutes{}
}

func (r *HealthcheckRoutes) RegisterRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/healthcheck", append(middleware.BaseMiddlewares, handler.Handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheck.NewHealthCheckHandler()))...)
}
