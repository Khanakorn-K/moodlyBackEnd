package models

import "time"

type MoodLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Mood      int       `gorm:"not null" json:"mood"`
	Note      string    `json:"note"`
	Causes    string    `json:"causes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
