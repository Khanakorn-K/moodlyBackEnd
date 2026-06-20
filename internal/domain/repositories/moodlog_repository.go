package repositories

import (
	"moodly/internal/domain/entities"
	"time"
)

type MoodLogRepositoryInterface interface {
	CreateMoodLog(moodLog *entities.MoodLogEntity) error

	FindMoodLogsByDate(userID uint, date time.Time) ([]entities.MoodLogEntity, error)

	UpdateMoodLog(moodLog *entities.MoodLogEntity) error

	DeleteMoodLog(id uint, userID uint) error
}
