package repositories

import (
	models "moodly/Models"

	"gorm.io/gorm"
)

type OverviewRepository struct {
	db *gorm.DB
}

func NewOverviewRepository(db *gorm.DB) *OverviewRepository {
	return &OverviewRepository{db: db}
}

func (r *OverviewRepository) FindMoodLogsByDateRange(
	userID uint,
	startDate string,
	endDate string,
) ([]models.MoodLog, error) {
	var moodLogs []models.MoodLog

	err := r.db.
		Where("user_id = ?", userID).
		Where("created_at >= ?", startDate+" 00:00:00").
		Where("created_at <= ?", endDate+" 23:59:59").
		Order("created_at ASC").
		Find(&moodLogs).Error

	if err != nil {
		return nil, err
	}

	return moodLogs, nil
}
