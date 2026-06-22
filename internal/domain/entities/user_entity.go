package entities

import "time"

type UserEntity struct {
	ID           uint                 `gorm:"primaryKey" json:"id"`
	Name         string               `gorm:"type:varchar(100);not null" json:"name"`
	Email        string               `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password     *string              `gorm:"type:varchar(255)" json:"-"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	MoodLogs     []MoodLogEntity      `gorm:"foreignKey:UserID" json:"mood_logs,omitempty"`
	CustomCauses []CustomCauseEntity  `gorm:"foreignKey:UserID" json:"custom_causes,omitempty"`
	Accounts     []OAuthAccountEntity `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
}
