package models

import "time"

type ScanSession struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	FamilyID  string    `json:"family_id" gorm:"column:family_id"`
	Status	string    `json:"status" gorm:"column:status"`
	TargetMembers int       `json:"target_members" gorm:"column:target_members"`
	CurrentUploads int       `json:"current_uploads" gorm:"column:current_uploads"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`	
}


func (ScanSession) TableName() string {
	return "scan_sessions"
}