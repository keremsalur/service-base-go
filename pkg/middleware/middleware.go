package middleware

import (
	"github.com/gofiber/fiber/v2"
)

var BaseMiddlewares = []fiber.Handler{
	RequestLogger,
	//JwtMiddleware,
	DynamicDTOValidationMiddleware,
}
