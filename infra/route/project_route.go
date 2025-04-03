package route

import (
	"service-base-go/app/project"
	"service-base-go/infra/repository"
	"service-base-go/pkg/handler"
	"service-base-go/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProjectRoutes struct {
}

func NewProjectRoutes() Routes {
	return &ProjectRoutes{}
}

func (r *ProjectRoutes) RegisterRoutes(app *fiber.App, db *gorm.DB) {
	projectRoutesV1 := app.Group("/api/v1/project")
	projectRoutesV1.Post("/", append(middleware.BaseMiddlewares, handler.Handle[project.CreateProjectRequest, project.CreateProjectResponse](project.NewCreateProductHandler(repository.NewProjectRepository(db))))...)

}
