package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"strings"
)

type MoodLogsService struct {
	repo *repositories.MoodLogRepository
}

func NewMoodLogsService(repo *repositories.MoodLogRepository) *MoodLogsService {
	return &MoodLogsService{repo: repo}
}

func (s *MoodLogsService) CreateMoodLog(moodLog *models.MoodLog) error {
	moodLog.Note = strings.TrimSpace(moodLog.Note)
	moodLog.Causes = strings.TrimSpace(moodLog.Causes)

	if moodLog.UserID == 0 {
		return errors.New("user id is required")
	}

	if moodLog.Mood == 0 {
		return errors.New("mood is required")
	}

	return s.repo.CreateMoodLog(moodLog)
}

func (s *MoodLogsService) GetMoodLogsByDate(userID uint, date string) ([]models.MoodLog, error) {
	date = strings.TrimSpace(date)

	if userID == 0 {
		return nil, errors.New("user id is required")
	}

	if date == "" {
		return nil, errors.New("date is required")
	}

	return s.repo.FindMoodLogsByDate(userID, date)
}

func (s *MoodLogsService) UpdateMoodLog(moodLog *models.MoodLog) error {
	if moodLog.ID == 0 {
		return errors.New("mood log id is required")
	}

	if moodLog.UserID == 0 {
		return errors.New("user id is required")
	}

	if moodLog.Mood == 0 {
		return errors.New("mood is required")
	}

	return s.repo.UpdateMoodLog(moodLog)
}

func (s *MoodLogsService) DeleteMoodLog(id uint, userID uint) error {
	if id == 0 {
		return errors.New("mood log id is required")
	}

	if userID == 0 {
		return errors.New("user id is required")
	}

	return s.repo.DeleteMoodLog(id, userID)
}
