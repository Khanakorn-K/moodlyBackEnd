package repositories

import (
	"errors"
	models "moodly/Models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) Login(email string, password string) (*models.User, error) {
	user, err := r.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) FindOAuthAccount(provider string, providerAccountID string) (*models.OAuthAccount, error) {
	var account models.OAuthAccount

	err := r.db.Where(
		"provider = ? AND provider_account_id = ?",
		provider,
		providerAccountID,
	).First(&account).Error

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AuthRepository) FindOAuthAccountWithUser(provider string, providerAccountID string) (*models.OAuthAccount, error) {
	var account models.OAuthAccount

	err := r.db.
		Preload("User").
		Where("provider = ? AND provider_account_id = ?", provider, providerAccountID).
		First(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AuthRepository) CreateOAuthAccount(account *models.OAuthAccount) error {
	return r.db.Create(account).Error
}
func (r *AuthRepository) FindOrCreateOAuthAccount(
	userID uint,
	provider string,
	providerAccountID string,
) (*models.OAuthAccount, error) {
	var account models.OAuthAccount

	err := r.db.
		Where("provider = ? AND provider_account_id = ?", provider, providerAccountID).
		First(&account).Error

	if err == nil {
		if account.UserID != userID {
			return nil, err
		}

		return &account, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	account = models.OAuthAccount{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
	}

	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
