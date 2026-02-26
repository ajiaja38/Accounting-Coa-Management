package report

import (
	"fiber.com/session-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, handler *Handler) {
	reportRoutes := router.Group("/report")
	reportRoutes.Use(middleware.AuthMiddleware())

	reportRoutes.Get("/ledger", handler.GetLedger)
	reportRoutes.Get("/trial-balance", handler.GetTrialBalance)
	reportRoutes.Get("/profit-loss", handler.GetProfitLoss)
	reportRoutes.Get("/balance-sheet", handler.GetBalanceSheet)
}
