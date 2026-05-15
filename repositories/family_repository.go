package repositories

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"nutapp-backend/database"
	"nutapp-backend/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateFamily(family *models.Family) error {
	return database.DB.Create(family).Error
}

func GenerateUniqueFamilyCode() (string, error) {
	const minLength = 6
	const maxLength = 8
	const maxAttempts = 25
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	for attempt := 0; attempt < maxAttempts; attempt++ {
		length, err := randomIntInRange(minLength, maxLength)
		if err != nil {
			return "", err
		}

		code, err := generateCode(alphabet, length)
		if err != nil {
			return "", err
		}

		exists, err := familyCodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", errors.New("no se pudo generar un código de familia único")
}

func ProfileExistsByID(profileID string) (bool, error) {
	trimmedProfileID := strings.TrimSpace(profileID)
	if trimmedProfileID == "" {
		return false, errors.New("profile_id es obligatorio")
	}

	var count int64
	if err := database.DB.Model(&models.Profile{}).Where("id = ?", trimmedProfileID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("error verificando profile autenticado: %w", err)
	}

	return count > 0, nil
}

func familyCodeExists(code string) (bool, error) {
	var count int64
	if err := database.DB.Model(&models.Family{}).Where("family_code = ?", code).Count(&count).Error; err != nil {
		return false, fmt.Errorf("error verificando unicidad del código: %w", err)
	}

	return count > 0, nil
}

func randomIntInRange(min, max int) (int, error) {
	if min > max {
		return 0, errors.New("rango inválido para longitud del código")
	}

	rangeSize := max - min + 1
	value, err := rand.Int(rand.Reader, big.NewInt(int64(rangeSize)))
	if err != nil {
		return 0, fmt.Errorf("error generando longitud aleatoria: %w", err)
	}

	return min + int(value.Int64()), nil
}

func generateCode(alphabet string, length int) (string, error) {
	if length <= 0 {
		return "", errors.New("la longitud del código debe ser mayor a cero")
	}

	result := make([]byte, length)
	for index := 0; index < length; index++ {
		value, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", fmt.Errorf("error generando carácter aleatorio: %w", err)
		}

		result[index] = alphabet[value.Int64()]
	}

	return string(result), nil
}

func EnsureProfileFamily(profileID, familyName string) (*models.Family, bool, error) {
	trimmedProfileID := strings.TrimSpace(profileID)
	if trimmedProfileID == "" {
		return nil, false, errors.New("profile_id es obligatorio")
	}

	trimmedFamilyName := strings.TrimSpace(familyName)
	var family models.Family
	created := false

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var profile models.Profile
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", trimmedProfileID).First(&profile).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("no se encontro profile para el usuario %s", trimmedProfileID)
			}
			return fmt.Errorf("error buscando profile para asignar familia: %w", err)
		}

		if profile.FamilyID != nil && strings.TrimSpace(*profile.FamilyID) != "" {
			if err := tx.Where("id = ?", *profile.FamilyID).First(&family).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("no se encontro la familia asociada al profile %s", trimmedProfileID)
				}
				return fmt.Errorf("error buscando familia existente: %w", err)
			}
			return nil
		}

		if trimmedFamilyName == "" {
			trimmedFamilyName = "Familia"
		}

		generatedFamily, err := createFamilyWithUniqueCode(tx, trimmedFamilyName)
		if err != nil {
			return err
		}

		profile.FamilyID = &generatedFamily.ID
		if err := tx.Save(&profile).Error; err != nil {
			return fmt.Errorf("error asignando familia al profile: %w", err)
		}

		family = *generatedFamily
		created = true
		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return &family, created, nil
}

func createFamilyWithUniqueCode(tx *gorm.DB, familyName string) (*models.Family, error) {
	const maxAttempts = 25

	for attempt := 0; attempt < maxAttempts; attempt++ {
		code, err := GenerateUniqueFamilyCode()
		if err != nil {
			return nil, err
		}

		family := models.Family{
			ID:         uuid.NewString(),
			Name:       familyName,
			FamilyCode: code,
		}

		if err := tx.Create(&family).Error; err != nil {
			if isUniqueConstraintError(err) {
				continue
			}

			return nil, fmt.Errorf("error creando familia: %w", err)
		}

		return &family, nil
	}

	return nil, errors.New("no se pudo crear una familia con un código único")
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key value violates unique constraint") || strings.Contains(message, "sqlstate 23505")
}
