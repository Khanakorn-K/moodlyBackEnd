package repositories

import (
	"moodly/internal/domain/entities"
	"time"
)

type InsightRepositoryInterface interface {
	FindInsightLogs(userID uint, selectedDate *time.Time) ([]entities.MoodLogEntity, error)
}
