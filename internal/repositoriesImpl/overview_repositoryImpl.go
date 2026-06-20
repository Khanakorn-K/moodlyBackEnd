package repositoriesImpl

import (
	"moodly/internal/domain/entities"
	"moodly/internal/domain/repositories"
	"time"

	"gorm.io/gorm"
)

type OverviewRepositoryImpl struct {
	db *gorm.DB
}

func NewOverviewRepositoryImpl(db *gorm.DB) repositories.OverviewRepositoryInterface {
	return &OverviewRepositoryImpl{db: db}
}

func (r *OverviewRepositoryImpl) FindMoodLogsByDateRange(
	userID uint,
	startDate time.Time,
	endDate time.Time,
) ([]entities.MoodLogEntity, error) {
	var moodLogs []entities.MoodLogEntity

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
