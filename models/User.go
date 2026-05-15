package models

import "time"

/*
	This will be the model for the supabase auth user, we will use the UID as the primary key and we will store the name and email in our own database, we will not store the password in our database, we will use the supabase auth to handle the authentication and authorization of the users, we will only store the UID, name and email in our database for reference and to link with other models like family, meals, etc.
*/
type User struct {
	ID              string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Email           string    `json:"email" gorm:"column:email"`
	EmailConfirmedAt *time.Time `json:"email_confirmed_at" gorm:"column:email_confirmed_at"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at"`
	
}


func (User) TableName() string {
	return "auth.users"
}
