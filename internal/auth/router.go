package auth

import (
	"fiber.com/session-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, handler *Handler) {
	auth := router.Group("/auth")

	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	auth.Use(middleware.AuthMiddleware())
	auth.Post("/logout", handler.Logout)
	auth.Get("/me", handler.Me)
}
