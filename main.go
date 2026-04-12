package main

import (
	"nutapp-backend/database"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "log"
)

func main() {
    database.ConnectDB()
    app := fiber.New()
    app.Use(cors.New())
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "online",
            "message": "Motor Go encendido 🚀",
        })
    })
    log.Fatal(app.Listen(":8080"))
}