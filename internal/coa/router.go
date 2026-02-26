package coa

import (
	"fiber.com/session-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, handler *Handler) {
	coaRoutes := router.Group("/coa")
	coaRoutes.Use(middleware.AuthMiddleware())

	coaRoutes.Get("/", handler.GetAll)
	coaRoutes.Get("/with-children", handler.GetAllWithChildren)
	coaRoutes.Get("/:code", handler.GetByCode)
	coaRoutes.Post("/", handler.Create)
	coaRoutes.Put("/:code", handler.Update)
	coaRoutes.Delete("/:code", handler.Delete)
}
