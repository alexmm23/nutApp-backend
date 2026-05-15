package handlers

import (
	"nutapp-backend/repositories"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GoogleProfileRequest struct {
	GoogleID  string `json:"google_id"`
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "name, email y password son obligatorios",
		})
	}

	user, err := repositories.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Usuario creado exitosamente",
		"data":    user,
	})

}

func UpsertGoogleProfile(c *fiber.Ctx) error {
	println("Recibiendo request para upsert Google Profile")
	var req GoogleProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	req.GoogleID = strings.TrimSpace(req.GoogleID)
	req.Name = strings.TrimSpace(req.Name)
	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = strings.TrimSpace(req.Email)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)

	if req.Name == "" {
		req.Name = req.FullName
	}

	if req.Name == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "name y email son obligatorios",
		})
	}
	println("Request validado, llamando a repositorio para upsert Google Profile")

	profile, created, err := repositories.UpsertGoogleProfile(req.GoogleID, req.Name, req.Email, req.AvatarURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	statusCode := 200
	message := "Perfil sincronizado exitosamente"
	if created {
		statusCode = 201
		message = "Perfil creado exitosamente"
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"message":    message,
		"created":    created,
		"profile_id": profile.ID,
		"family_id":  profile.FamilyID,
		"data":       profile,
	})
}

func GetFamilyMembersByUserID(c *fiber.Ctx) error {
	userID := strings.TrimSpace(c.Params("user_id"))
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "user_id es obligatorio",
		})
	}

	members, err := repositories.GetFamilyMembersByUserID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "no se encontro profile") {
			return c.Status(404).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Miembros de familia obtenidos exitosamente",
		"data":    members,
	})
}
