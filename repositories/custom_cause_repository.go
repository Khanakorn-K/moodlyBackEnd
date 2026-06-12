package repositories

import (
	"errors"
	models "moodly/Models"

	"gorm.io/gorm"
)

type CustomCauseRepository struct {
	db *gorm.DB
}

func NewCustomCauseRepository(db *gorm.DB) *CustomCauseRepository {
	return &CustomCauseRepository{db: db}
}

func (r *CustomCauseRepository) Create(cause *models.CustomCause) error {
	return r.db.Create(cause).Error
}

func (r *CustomCauseRepository) FindByUserID(userID uint) ([]models.CustomCause, error) {
	var causes []models.CustomCause

	err := r.db.Where("user_id = ?", userID).Find(&causes).Error
	if err != nil {
		return nil, err
	}

	return causes, nil
}

func (r *CustomCauseRepository) FindByID(id uint, userID uint) (*models.CustomCause, error) {
	var cause models.CustomCause

	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&cause).Error
	if err != nil {
		return nil, err
	}

	return &cause, nil
}

func (r *CustomCauseRepository) Update(cause *models.CustomCause) error {
	result := r.db.
		Model(&models.CustomCause{}).
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

func (r *CustomCauseRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CustomCause{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("cause not found or unauthorized")
	}

	return nil
}
