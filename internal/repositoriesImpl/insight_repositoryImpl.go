package repositoriesImpl

import (
	"moodly/internal/domain/entities"
	"moodly/internal/domain/repositories"
	"time"

	"gorm.io/gorm"
)

type InsightRepositoryImpl struct {
	db *gorm.DB
}

func NewInsightRepositoryImpl(db *gorm.DB) repositories.InsightRepositoryInterface {
	return &InsightRepositoryImpl{db: db}
}

func (r *InsightRepositoryImpl) FindInsightLogs(
	userID uint,
	selectedDate *time.Time,
) ([]entities.MoodLogEntity, error) {
	var moodLogs []entities.MoodLogEntity

	query := r.db.Where("user_id = ?", userID)

	if selectedDate != nil {
		startDate := *selectedDate
		endDate := startDate.AddDate(0, 0, 1)

		query = query.Where("created_at >= ? AND created_at < ?", startDate, endDate)
	}

	if err := query.Find(&moodLogs).Error; err != nil {
		return nil, err
	}

	return moodLogs, nil
}
