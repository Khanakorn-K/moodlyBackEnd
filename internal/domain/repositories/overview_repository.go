package repositories

import (
	"moodly/internal/domain/entities"
	"time"
)

type OverviewRepositoryInterface interface {
	FindMoodLogsByDateRange(userID uint, start time.Time, end time.Time) ([]entities.MoodLogEntity, error)
}
