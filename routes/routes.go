package routes

import (
	"nutapp-backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.Health)
	app.Get("/debug/db", handlers.DebugDB)
	app.Post("/new/family", handlers.CreateFamily)
}
