package models

import "time"

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password     *string        `gorm:"type:varchar(255)" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	MoodLogs     []MoodLog      `gorm:"foreignKey:UserID" json:"mood_logs,omitempty"`
	CustomCauses []CustomCause  `gorm:"foreignKey:UserID" json:"custom_causes,omitempty"`
	Accounts     []OAuthAccount `gorm:"foreignKey:UserID" json:"accounts,omitempty"`
}
type OAuthAccount struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`

	Provider          string `gorm:"type:varchar(50);not null;uniqueIndex:idx_oauth_provider_account" json:"provider"`
	ProviderAccountID string `gorm:"type:varchar(255);not null;uniqueIndex:idx_oauth_provider_account" json:"provider_account_id"`

	AccessToken  string     `gorm:"type:text" json:"-"`
	RefreshToken string     `gorm:"type:text" json:"-"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
