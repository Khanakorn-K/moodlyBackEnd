package repositoriesImpl

import (
	"errors"
	"moodly/internal/domain/entities"
	"moodly/internal/domain/repositories"

	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryImpl(db *gorm.DB) repositories.AuthRepositoryInterface {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) FindByEmail(email string) (*entities.UserEntity, error) {
	var user entities.UserEntity

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) CreateUser(user *entities.UserEntity) error {
	return r.db.Create(user).Error
}

func (r *AuthRepositoryImpl) FindOrCreateOAuthAccount(
	userID uint,
	provider string,
	providerAccountID string,
) (*entities.OAuthAccountEntity, error) {
	var account entities.OAuthAccountEntity

	err := r.db.
		Where("provider = ? AND provider_account_id = ?", provider, providerAccountID).
		First(&account).Error

	if err == nil {
		if account.UserID != userID {
			return nil, errors.New("oauth account already linked to another user")
		}

		return &account, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	account = entities.OAuthAccountEntity{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
	}

	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
