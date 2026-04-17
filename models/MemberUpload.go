package models

import "time"

type MemberUpload struct {
	ID int `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
	SessionID string    `json:"session_id" gorm:"column:session_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	ImageURL string    `json:"image_url" gorm:"column:image_url"`
	RawOCRData string    `json:"raw_ocr_data" gorm:"column:raw_ocr_data"`
}


func (MemberUpload) TableName() string {
	return "member_uploads"
}
