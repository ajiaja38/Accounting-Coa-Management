package journal

import (
	"fiber.com/session-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, handler *Handler, db *gorm.DB) {
	journalRoutes := router.Group("/journal")
	journalRoutes.Use(middleware.AuthMiddleware())

	journalRoutes.Get("/", handler.GetAll)
	journalRoutes.Get("/:id", handler.GetByID)
	journalRoutes.Put("/:id/post", handler.PostJournal)
	journalRoutes.Delete("/:id", handler.Delete)

	journalRoutes.Post("/", middleware.DBTransaction(db), handler.Create)
}
