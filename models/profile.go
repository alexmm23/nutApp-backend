package models

import "time"

type Profile struct {
	ID		string    `json:"id" gorm:"column:id;primaryKey"`
	FullName string    `json:"full_name" gorm:"column:full_name"`
	Email   string    `json:"email" gorm:"column:email;unique"`
	AvatarURL string    `json:"avatar_url" gorm:"column:avatar_url"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	FamilyID string    `json:"family_id" gorm:"column:family_id"`



}