package handlers

import (
	"log"
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

type GenerateFamilyCodeRequest struct {
	ProfileID string  `json:"profile_id"`
	FamilyID  *string `json:"family_id"`
	GoogleID  string  `json:"google_id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	AvatarURL string  `json:"avatar_url"`
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

func GenerateFamilyCode(c *fiber.Ctx) error {
	var req GenerateFamilyCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	req.ProfileID = strings.TrimSpace(req.ProfileID)
	req.GoogleID = strings.TrimSpace(req.GoogleID)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)

	profileID := req.ProfileID
	if profileID == "" {
		authorization := strings.TrimSpace(c.Get("Authorization"))
		if strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
			profileID = strings.TrimSpace(authorization[7:])
		}
	}

	if profileID == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "autenticación requerida",
		})
	}

	registered, err := repositories.ProfileExistsByID(profileID)
	if err != nil {
		log.Println("error validando autenticación de perfil:", err)
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !registered {
		return c.Status(403).JSON(fiber.Map{
			"error": "el perfil no está registrado en la base de datos",
		})
	}

	if req.ProfileID == "" {
		req.ProfileID = profileID
	}

	code, err := repositories.GenerateUniqueFamilyCode()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":     "Código de familia generado exitosamente",
		"family_code": code,
		"profile_id":  profileID,
	})
}
