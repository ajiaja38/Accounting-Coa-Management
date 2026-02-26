package main

import (
	"fmt"
	"log"

	"fiber.com/session-api/config"
	_ "fiber.com/session-api/docs"
	"fiber.com/session-api/internal/auth"
	"fiber.com/session-api/internal/coa"
	"fiber.com/session-api/internal/domain"
	"fiber.com/session-api/internal/journal"
	"fiber.com/session-api/internal/report"
	"fiber.com/session-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title           Financial Accounting API
// @version         1.0
// @description     RESTful API for financial accounting management (COA-based) built with Go Fiber.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@fiber-coa.local

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
func main() {
	config.Load()

	config.ConnectDatabase()
	db := config.DB

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.ChartOfAccount{},
		&domain.JournalEntry{},
		&domain.JournalEntryDetail{},
	); err != nil {
		log.Fatalf("Auto-migrate failed: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:5173,http://localhost:8080",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api/v1")

	// Auth routes
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(api, authHandler)

	// COA routes
	coaRepo := coa.NewRepository(db)
	coaService := coa.NewService(coaRepo)
	coaHandler := coa.NewHandler(coaService)
	coa.RegisterRoutes(api, coaHandler)

	// Journal routes
	journalRepo := journal.NewRepository(db)
	journalService := journal.NewService(journalRepo)
	journalHandler := journal.NewHandler(journalService)
	journal.RegisterRoutes(api, journalHandler, db)

	// Report routes
	reportRepo := report.NewRepository(db)
	reportService := report.NewService(reportRepo, coaRepo)
	reportHandler := report.NewHandler(reportService)
	report.RegisterRoutes(api, reportHandler)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Financial Accounting API is running",
		})
	})

	addr := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Swagger UI: http://localhost%s/swagger/index.html", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
