package handlers

import (
	"strings"

	"nutapp-backend/models"
	"nutapp-backend/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateFamilyRequest struct {
	FamilyName string `json:"family_name"`
	FamilyCode string `json:"family_code"`
}

func CreateFamily(c *fiber.Ctx) error {
	var req CreateFamilyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	req.FamilyName = strings.TrimSpace(req.FamilyName)
	req.FamilyCode = strings.TrimSpace(req.FamilyCode)

	if req.FamilyName == "" || req.FamilyCode == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "family_name y family_code son obligatorios",
		})
	}

	family := models.Family{
		ID:         uuid.NewString(),
		Name:       req.FamilyName,
		FamilyCode: req.FamilyCode,
	}

	if err := repositories.CreateFamily(&family); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Familia creada exitosamente",
		"data":    family,
	})
}
