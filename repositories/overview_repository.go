package repositories

import (
	models "moodly/Models"
	"time"

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
	startDate time.Time,
	endDate time.Time,
) ([]models.MoodLog, error) {
	var moodLogs []models.MoodLog

	err := r.db.
		Where("user_id = ?", userID).
		Where("created_at >= ?", startDate).
		Where("created_at < ?", endDate.AddDate(0, 0, 1)).
		Order("created_at ASC").
		Find(&moodLogs).Error

	if err != nil {
		return nil, err
	}

	return moodLogs, nil
}
