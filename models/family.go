package models

import "time"

type Family struct {
	ID         string    `json:"id" gorm:"column:id;primaryKey"`
	Name       string    `json:"family_name" gorm:"column:family_name"`
	FamilyCode string    `json:"family_code" gorm:"column:family_code"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

func (Family) TableName() string {
	return "families"
}