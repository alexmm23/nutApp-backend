package models

import "time"


type WeeklyDiet struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
	DayOfWeek string    `json:"day_of_week" gorm:"column:day_of_week"`
	MealTime  string    `json:"meal_time" gorm:"column:meal_time"`
	FoodName  string    `json:"food_name" gorm:"column:food_name"`
	Calories  int       `json:"calories" gorm:"column:kcal"`
	Description string    `json:"description" gorm:"column:description"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	FamilyID  string    `json:"family_id" gorm:"column:family_id"`
}