package models

import "time"

type User struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"type:varchar(100);not null" json:"name"`
	Email        string        `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password     string        `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	MoodLogs     []MoodLog     `gorm:"foreignKey:UserID" json:"mood_logs,omitempty"`
	CustomCauses []CustomCause `gorm:"foreignKey:UserID" json:"custom_causes,omitempty"`
}
