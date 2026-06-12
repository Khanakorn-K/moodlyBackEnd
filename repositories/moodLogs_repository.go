package repositories

import (
	"errors"
	models "moodly/Models"
	"time"

	"gorm.io/gorm"
)

type MoodLogRepository struct {
	db *gorm.DB
}

func NewMoodLogsRepository(db *gorm.DB) *MoodLogRepository {
	return &MoodLogRepository{db: db}
}

func (r *MoodLogRepository) CreateMoodLog(moodLog *models.MoodLog) error {
	return r.db.Create(moodLog).Error
}

func (r *MoodLogRepository) FindMoodLogsByDate(userID uint, date time.Time) ([]models.MoodLog, error) {
	var moodLogs []models.MoodLog

	err := r.db.
		Where("user_id = ?", userID).
		Where("created_at >= ?", date).
		Where("created_at < ?", date.AddDate(0, 0, 1)).
		Find(&moodLogs).Error

	if err != nil {
		return nil, err
	}

	return moodLogs, nil
}

func (r *MoodLogRepository) UpdateMoodLog(moodLog *models.MoodLog) error {
	result := r.db.
		Model(&models.MoodLog{}).
		Where("id = ? AND user_id = ?", moodLog.ID, moodLog.UserID).
		Updates(map[string]interface{}{
			"mood":   moodLog.Mood,
			"note":   moodLog.Note,
			"causes": moodLog.Causes,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("mood log not found or unauthorized")
	}

	return nil
}

func (r *MoodLogRepository) DeleteMoodLog(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.MoodLog{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("mood log not found or unauthorized")
	}

	return nil
}
