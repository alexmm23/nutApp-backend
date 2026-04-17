package repositories

import (
	"nutapp-backend/database"
	"nutapp-backend/models"
)

func CreateFamily(family *models.Family) error {
	return database.DB.Create(family).Error
}
