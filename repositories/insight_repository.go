package repositories

import (
	models "moodly/Models"

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
	mood string,
	startDate string,
	endDate string,
) ([]models.MoodLog, int64, error) {
	var moodLogs []models.MoodLog
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if mood != "" {
		query = query.Where("mood = ?", mood)
	}

	if startDate != "" {
		query = query.Where("created_at >= ?", startDate+" 00:00:00")
	}

	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	if err := query.Model(&models.MoodLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Find(&moodLogs).Error; err != nil {
		return nil, 0, err
	}

	return moodLogs, total, nil
}
