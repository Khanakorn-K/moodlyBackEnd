package repositories

import (
	models "moodly/Models"
	"time"

	"gorm.io/gorm"
)

type InsightRepository struct {
	db *gorm.DB
}

func NewInsightRepository(db *gorm.DB) *InsightRepository {
	return &InsightRepository{db: db}
}

func (r *InsightRepository) FindMoodLogs(
	userID uint,
	mood *int,
	startDate *time.Time,
	endDate *time.Time,
) ([]models.MoodLog, int64, error) {
	var moodLogs []models.MoodLog
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if mood != nil {
		query = query.Where("mood = ?", *mood)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at < ?", endDate.AddDate(0, 0, 1))
	}

	if err := query.Model(&models.MoodLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Find(&moodLogs).Error; err != nil {
		return nil, 0, err
	}

	return moodLogs, total, nil
}
