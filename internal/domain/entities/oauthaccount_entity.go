package entities

import "time"

type OAuthAccountEntity struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`

	Provider          string `gorm:"type:varchar(50);not null;uniqueIndex:idx_oauth_provider_account" json:"provider"`
	ProviderAccountID string `gorm:"type:varchar(255);not null;uniqueIndex:idx_oauth_provider_account" json:"provider_account_id"`

	AccessToken  string     `gorm:"type:text" json:"-"`
	RefreshToken string     `gorm:"type:text" json:"-"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User UserEntity `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
