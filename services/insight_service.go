package services

import (
	"encoding/json"
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidDateFormat = errors.New("invalid date format")

type InsightService struct {
	repo *repositories.InsightRepository
}

func NewInsightService(repo *repositories.InsightRepository) *InsightService {
	return &InsightService{repo: repo}
}

type InsightResult struct {
	TotalLogs        int                       `json:"totalLogs"`
	AverageMood      float64                   `json:"averageMood"`
	MoodDistribution map[string]int            `json:"moodDistribution"`
	CausesAnalysis   map[string]map[string]int `json:"causesAnalysis"`
}

func (s *InsightService) GetInsights(userID uint, selectedDate string) (*InsightResult, error) {
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

	return &InsightResult{
		TotalLogs:        len(logs),
		AverageMood:      calculateAverageMood(logs),
		MoodDistribution: calculateMoodDistributionRecord(logs),
		CausesAnalysis:   calculateCauseAnalysis(logs),
	}, nil
}

func calculateMoodDistributionRecord(logs []models.MoodLog) map[string]int {
	result := map[string]int{}

	for _, log := range logs {
		moodKey := strconv.Itoa(log.Mood)
		result[moodKey]++
	}

	return result
}

func calculateCauseAnalysis(logs []models.MoodLog) map[string]map[string]int {
	result := map[string]map[string]int{}

	for _, log := range logs {
		causes := parseCauses(log.Causes)
		moodKey := strconv.Itoa(log.Mood)

		for _, cause := range causes {
			cause = strings.TrimSpace(cause)
			if cause == "" {
				continue
			}

			if result[cause] == nil {
				result[cause] = map[string]int{}
			}

			result[cause][moodKey]++
		}
	}

	return result
}

func parseCauses(rawCauses string) []string {
	rawCauses = strings.TrimSpace(rawCauses)

	if rawCauses == "" {
		return []string{}
	}

	var causes []string

	// กรณีเก็บเป็น JSON string เช่น ["work","money"]
	if err := json.Unmarshal([]byte(rawCauses), &causes); err == nil {
		return causes
	}

	// กรณีเก็บเป็น string ธรรมดา เช่น work,money
	parts := strings.Split(rawCauses, ",")

	cleanCauses := []string{}

	for _, part := range parts {
		cause := strings.TrimSpace(part)
		cause = strings.Trim(cause, `"'[]{} `)

		if cause != "" {
			cleanCauses = append(cleanCauses, cause)
		}
	}

	return cleanCauses
}
