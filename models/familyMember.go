package models

import "time"

type FamilyMember struct {
	ID 	 string `json:"id" gorm:"column:id;primaryKey"`
	userID string `json:"user_id" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:name"`
	relationship string `json:"relationship" gorm:"column:relationship"`
	AvatarURL string `json:"avatar_url" gorm:"column:avatar_url"`
	createdAt time.Time `json:"created_at" gorm:"column:created_at"`


}