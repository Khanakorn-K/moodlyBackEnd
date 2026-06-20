package repositories

import "moodly/internal/domain/entities"

type CustomCauseRepositoryInterface interface {
	Create(cause *entities.CustomCauseEntity) error
	FindByUserID(userID uint) (*[]entities.CustomCauseEntity, error)
	Update(cause *entities.CustomCauseEntity) error
	Delete(id uint, userID uint) error
}
