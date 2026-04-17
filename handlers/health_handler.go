package handlers

import (
	"log"
	"nutapp-backend/database"

	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "online",
		"message": "Motor Go encendido 🚀",
	})
}

func DebugDB(c *fiber.Ctx) error {
	var now string

	if err := database.DB.Raw("SELECT NOW()::text").Scan(&now).Error; err != nil {
		log.Println("error consultando la base de datos:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Println("resultado de la consulta:", now)

	return c.JSON(fiber.Map{
		"query":  "SELECT NOW()::text",
		"result": now,
	})
}
