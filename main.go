package main

import (
	"log"
	"nutapp-backend/database"
	"nutapp-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDB()

	app := fiber.New()
	app.Use(cors.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}