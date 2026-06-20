package repositories

import "moodly/internal/domain/entities"

type AuthRepositoryInterface interface {
	FindByEmail(email string) (*entities.UserEntity, error)

	CreateUser(user *entities.UserEntity) error

	FindOrCreateOAuthAccount(
		userID uint,
		provider string,
		providerAccountID string,
	) (*entities.OAuthAccountEntity, error)
}
