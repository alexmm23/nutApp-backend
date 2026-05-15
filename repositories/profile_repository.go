package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"nutapp-backend/database"
	"nutapp-backend/models"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type supabaseCreateUserRequest struct {
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	EmailConfirm bool           `json:"email_confirm"`
	UserMetadata map[string]any `json:"user_metadata,omitempty"`
}

type supabaseCreateUserResponse struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Msg              string `json:"msg"`
	ErrorDescription string `json:"error_description"`
	User             struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func CreateUser(name, email, password string) (*models.User, error) {
	supabaseURL := strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
	serviceRoleKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || serviceRoleKey == "" {
		return nil, errors.New("faltan SUPABASE_URL o SUPABASE_SERVICE_ROLE_KEY en variables de entorno")
	}

	payload := supabaseCreateUserRequest{
		Email:        email,
		Password:     password,
		EmailConfirm: false,
		UserMetadata: map[string]any{"name": name},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload de supabase: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, supabaseURL+"/auth/v1/admin/users", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creando request a supabase: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("apikey", serviceRoleKey)
	request.Header.Set("Authorization", "Bearer "+serviceRoleKey)

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error llamando supabase auth: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de supabase: %w", err)
	}

	var parsedResponse supabaseCreateUserResponse
	if len(responseBody) > 0 {
		_ = json.Unmarshal(responseBody, &parsedResponse)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		message := strings.TrimSpace(parsedResponse.ErrorDescription)
		if message == "" {
			message = strings.TrimSpace(parsedResponse.Msg)
		}
		if message == "" {
			message = string(responseBody)
		}
		return nil, fmt.Errorf("supabase auth devolvio %d: %s", response.StatusCode, message)
	}

	userID := parsedResponse.User.ID
	if userID == "" {
		userID = parsedResponse.ID
	}

	userEmail := parsedResponse.User.Email
	if userEmail == "" {
		userEmail = parsedResponse.Email
	}
	if userEmail == "" {
		userEmail = email
	}

	if userID == "" {
		return nil, errors.New("supabase auth no devolvio id de usuario")
	}

	if _, _, err := upsertProfile(userID, name, userEmail, ""); err != nil {
		return nil, fmt.Errorf("usuario creado en auth, pero fallo guardando profile local: %w", err)
	}

	return &models.User{
		ID:    userID,
		Email: userEmail,
	}, nil
}

func UpsertGoogleProfile(googleID, name, email, avatarURL string) (*models.Profile, bool, error) {
	if email == "" {
		return nil, false, errors.New("email es obligatorio para crear o actualizar el profile")
	}

	if name == "" {
		return nil, false, errors.New("name es obligatorio para crear o actualizar el profile")
	}

	return upsertProfile(googleID, name, email, avatarURL)
}

func upsertProfile(preferredID, fullName, email, avatarURL string) (*models.Profile, bool, error) {
	trimmedEmail := strings.ToLower(strings.TrimSpace(email))
	trimmedName := strings.TrimSpace(fullName)
	trimmedAvatarURL := strings.TrimSpace(avatarURL)
	now := time.Now().UTC()

	if trimmedEmail == "" {
		return nil, false, errors.New("email es obligatorio")
	}

	var existing models.Profile
	err := database.DB.Where("email = ?", trimmedEmail).First(&existing).Error
	if err == nil {
		existing.FullName = trimmedName
		existing.Email = trimmedEmail
		existing.AvatarURL = trimmedAvatarURL
		existing.UpdatedAt = now

		if saveErr := database.DB.Save(&existing).Error; saveErr != nil {
			return nil, false, fmt.Errorf("error actualizando profile existente: %w", saveErr)
		}

		return &existing, false, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, fmt.Errorf("error buscando profile existente: %w", err)
	}

	profileID := strings.TrimSpace(preferredID)
	if profileID == "" {
		profileID = uuid.NewString()
	}

	profile := models.Profile{
		ID:        profileID,
		FullName:  trimmedName,
		Email:     trimmedEmail,
		AvatarURL: trimmedAvatarURL,
		UpdatedAt: now,
	}

	if createErr := database.DB.Create(&profile).Error; createErr != nil {
		return nil, false, fmt.Errorf("error creando profile: %w", createErr)
	}

	return &profile, true, nil
}
