package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"strings"
)

type InsightService struct {
	repo *repositories.InsightRepository
}

func NewInsightService(repo *repositories.InsightRepository) *InsightService {
	return &InsightService{repo: repo}
}

type MoodLogsResult struct {
	Data  []models.MoodLog `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
}

func (s *InsightService) FindMoodLogs(
	userID uint,
	mood string,
	startDate string,
	endDate string,
) (*MoodLogsResult, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}

	mood = strings.TrimSpace(mood)
	startDate = strings.TrimSpace(startDate)
	endDate = strings.TrimSpace(endDate)

	logs, total, err := s.repo.FindMoodLogs(userID, mood, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &MoodLogsResult{
		Data:  logs,
		Total: total,
		Page:  1,
	}, nil
}
