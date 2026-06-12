package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"strconv"
	"strings"
	"time"
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

	moodFilter, err := parseOptionalMood(mood)
	if err != nil {
		return nil, err
	}

	startDateFilter, err := parseOptionalDate(startDate)
	if err != nil {
		return nil, err
	}

	endDateFilter, err := parseOptionalDate(endDate)
	if err != nil {
		return nil, err
	}

	if startDateFilter != nil && endDateFilter != nil && startDateFilter.After(*endDateFilter) {
		return nil, errors.New("invalid date range")
	}

	logs, total, err := s.repo.FindMoodLogs(userID, moodFilter, startDateFilter, endDateFilter)
	if err != nil {
		return nil, err
	}

	return &MoodLogsResult{
		Data:  logs,
		Total: total,
		Page:  1,
	}, nil
}

func parseOptionalMood(value string) (*int, error) {
	if value == "" {
		return nil, nil
	}

	mood, err := strconv.Atoi(value)
	if err != nil || mood < 1 || mood > 5 {
		return nil, errors.New("mood must be between 1 and 5")
	}

	return &mood, nil
}

func parseOptionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	return &parsedDate, nil
}
