package route

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Routes interface {
	RegisterRoutes(app *fiber.App, db *gorm.DB)
}
