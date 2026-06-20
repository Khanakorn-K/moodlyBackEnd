package services

import (
	"errors"
	"moodly/internal/domain/entities"
	"moodly/internal/domain/repositories"
	"strings"
	"time"
)

var ErrInvalidDateFormat = errors.New("invalid date format")

type InsightService struct {
	repo repositories.InsightRepositoryInterface
}

func NewInsightService(repo repositories.InsightRepositoryInterface) *InsightService {
	return &InsightService{repo: repo}
}

func (s *InsightService) GetInsights(userID uint, selectedDate string) (*[]entities.MoodLogEntity, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}

	selectedDate = strings.TrimSpace(selectedDate)

	var selectedDateFilter *time.Time

	if selectedDate != "" {
		parsedDate, err := time.Parse("2006-01-02", selectedDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}

		selectedDateFilter = &parsedDate
	}

	logs, err := s.repo.FindInsightLogs(userID, selectedDateFilter)
	if err != nil {
		return nil, err
	}

	return &logs, nil
}
