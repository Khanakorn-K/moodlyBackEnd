package repositoriesImpl

import (
	"errors"
	"moodly/internal/domain/entities"
	"moodly/internal/domain/repositories"

	"gorm.io/gorm"
)

type CustomCauseRepositoryImpl struct {
	db *gorm.DB
}

func NewCustomCauseRepositoryImpl(db *gorm.DB) repositories.CustomCauseRepositoryInterface {
	return &CustomCauseRepositoryImpl{db: db}
}

func (r *CustomCauseRepositoryImpl) Create(cause *entities.CustomCauseEntity) error {
	return r.db.Create(cause).Error
}

func (r *CustomCauseRepositoryImpl) FindByUserID(userID uint) (*[]entities.CustomCauseEntity, error) {
	var causes []entities.CustomCauseEntity

	err := r.db.Where("user_id = ?", userID).Find(&causes).Error
	if err != nil {
		return nil, err
	}

	return &causes, nil
}

func (r *CustomCauseRepositoryImpl) FindByID(id uint, userID uint) (*entities.CustomCauseEntity, error) {
	var cause entities.CustomCauseEntity

	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&cause).Error
	if err != nil {
		return nil, err
	}

	return &cause, nil
}

func (r *CustomCauseRepositoryImpl) Update(cause *entities.CustomCauseEntity) error {
	result := r.db.
		Model(&entities.CustomCauseEntity{}).
		Where("id = ? AND user_id = ?", cause.ID, cause.UserID).
		Updates(map[string]interface{}{
			"name": cause.Name,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("cause not found or unauthorized")
	}

	return nil
}

func (r *CustomCauseRepositoryImpl) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entities.CustomCauseEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("cause not found or unauthorized")
	}

	return nil
}
