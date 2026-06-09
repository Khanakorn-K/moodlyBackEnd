package repositories

import (
	models "moodly/Models"

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

func (r *MoodLogRepository) FindMoodLogsByDate(userID uint, date string) ([]models.MoodLog, error) {
	var moodLogs []models.MoodLog

	err := r.db.
		Where("user_id = ? AND DATE(created_at) = ?", userID, date).
		Find(&moodLogs).Error

	if err != nil {
		return nil, err
	}

	return moodLogs, nil
}

func (r *MoodLogRepository) UpdateMoodLog(moodLog *models.MoodLog) error {
	return r.db.Save(moodLog).Error
}

func (r *MoodLogRepository) DeleteMoodLog(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.MoodLog{}).Error
}
