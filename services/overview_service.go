package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OverviewService struct {
	repo *repositories.OverviewRepository
}

func NewOverviewService(repo *repositories.OverviewRepository) *OverviewService {
	return &OverviewService{repo: repo}
}

type DailyCauseDistributionItem struct {
	Cause string `json:"cause"`
	Count int    `json:"count"`
}

type DailyMoodAverage struct {
	Date              string                       `json:"date"`
	AverageMood       float64                      `json:"averageMood"`
	TotalLogs         int                          `json:"totalLogs"`
	CauseDistribution []DailyCauseDistributionItem `json:"causeDistribution"`
}

type MoodDistributionItem struct {
	Mood  int `json:"mood"`
	Count int `json:"count"`
}

type MoodNote struct {
	ID        uint     `json:"id"`
	Date      string   `json:"date"`
	Mood      int      `json:"mood"`
	Note      string   `json:"note"`
	Causes    []string `json:"causes"`
	CreatedAt string   `json:"createdAt"`
}

type CauseSummary struct {
	Cause         string                 `json:"cause"`
	TotalCount    int                    `json:"totalCount"`
	MoodBreakdown []MoodDistributionItem `json:"moodBreakdown"`
}

type MonthlyAverageMoodResult struct {
	Month             string             `json:"month"`
	StartDate         string             `json:"startDate"`
	EndDate           string             `json:"endDate"`
	TotalLogs         int                `json:"totalLogs"`
	AverageMood       float64            `json:"averageMood"`
	DailyMoodAverages []DailyMoodAverage `json:"dailyMoodAverages"`
}

type OverviewResult struct {
	StartDate         string                 `json:"startDate"`
	EndDate           string                 `json:"endDate"`
	TotalLogs         int                    `json:"totalLogs"`
	AverageMood       float64                `json:"averageMood"`
	DailyMoodAverages []DailyMoodAverage     `json:"dailyMoodAverages"`
	MoodDistribution  []MoodDistributionItem `json:"moodDistribution"`
	MoodNotes         []MoodNote             `json:"moodNotes"`
	CauseSummaries    []CauseSummary         `json:"causeSummaries"`
}

func (s *OverviewService) GetMonthlyAverageMood(userID uint, month string) (*MonthlyAverageMoodResult, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}

	month = strings.TrimSpace(month)
	if month == "" {
		return nil, errors.New("month is required")
	}

	startDate, endDate, err := createMonthDateRange(month)
	if err != nil {
		return nil, err
	}

	logs, err := s.repo.FindMoodLogsByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	dateRange := createDateRange(startDate, endDate)

	return &MonthlyAverageMoodResult{
		Month:             month,
		StartDate:         formatDate(startDate),
		EndDate:           formatDate(endDate),
		TotalLogs:         len(logs),
		AverageMood:       calculateAverageMood(logs),
		DailyMoodAverages: calculateDailyMoodAverages(logs, dateRange),
	}, nil
}

func (s *OverviewService) GetOverview(userID uint, startDate string, endDate string) (*OverviewResult, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}

	startDate = strings.TrimSpace(startDate)
	endDate = strings.TrimSpace(endDate)
	if startDate == "" || endDate == "" {
		return nil, errors.New("date range is required")
	}

	parsedStartDate, err := parseDate(startDate)
	if err != nil {
		return nil, err
	}

	parsedEndDate, err := parseDate(endDate)
	if err != nil {
		return nil, err
	}

	if parsedStartDate.After(parsedEndDate) {
		return nil, errors.New("invalid date range")
	}

	logs, err := s.repo.FindMoodLogsByDateRange(userID, parsedStartDate, parsedEndDate)
	if err != nil {
		return nil, err
	}

	dateRange := createDateRange(parsedStartDate, parsedEndDate)

	return &OverviewResult{
		StartDate:         formatDate(parsedStartDate),
		EndDate:           formatDate(parsedEndDate),
		TotalLogs:         len(logs),
		AverageMood:       calculateAverageMood(logs),
		DailyMoodAverages: calculateDailyMoodAverages(logs, dateRange),
		MoodDistribution:  calculateMoodDistribution(logs),
		MoodNotes:         createMoodNotes(logs),
		CauseSummaries:    calculateCauseSummaries(logs),
	}, nil
}

func calculateAverageMood(logs []models.MoodLog) float64 {
	if len(logs) == 0 {
		return 0
	}

	total := 0
	for _, log := range logs {
		total += log.Mood
	}

	return roundOneDecimal(float64(total) / float64(len(logs)))
}

func calculateDailyMoodAverages(logs []models.MoodLog, dateRange []string) []DailyMoodAverage {
	type dailyValue struct {
		TotalMood   int
		TotalLogs   int
		CauseCounts map[string]int
	}

	dailyMap := map[string]*dailyValue{}

	for _, log := range logs {
		date := formatDate(log.CreatedAt)

		if dailyMap[date] == nil {
			dailyMap[date] = &dailyValue{
				TotalMood:   0,
				TotalLogs:   0,
				CauseCounts: map[string]int{},
			}
		}

		dailyMap[date].TotalMood += log.Mood
		dailyMap[date].TotalLogs++

		for _, cause := range splitCauses(log.Causes) {
			dailyMap[date].CauseCounts[cause]++
		}
	}

	result := []DailyMoodAverage{}

	for _, date := range dateRange {
		value := dailyMap[date]

		if value == nil {
			result = append(result, DailyMoodAverage{
				Date:              date,
				AverageMood:       0,
				TotalLogs:         0,
				CauseDistribution: []DailyCauseDistributionItem{},
			})
			continue
		}

		causeDistribution := []DailyCauseDistributionItem{}
		for cause, count := range value.CauseCounts {
			causeDistribution = append(causeDistribution, DailyCauseDistributionItem{
				Cause: cause,
				Count: count,
			})
		}

		sort.Slice(causeDistribution, func(i, j int) bool {
			return causeDistribution[i].Count > causeDistribution[j].Count
		})

		result = append(result, DailyMoodAverage{
			Date:              date,
			AverageMood:       roundOneDecimal(float64(value.TotalMood) / float64(value.TotalLogs)),
			TotalLogs:         value.TotalLogs,
			CauseDistribution: causeDistribution,
		})
	}

	return result
}

func calculateMoodDistribution(logs []models.MoodLog) []MoodDistributionItem {
	moodValues := []int{1, 2, 3, 4, 5}
	counts := map[int]int{}

	for _, log := range logs {
		counts[log.Mood]++
	}

	result := []MoodDistributionItem{}

	for _, mood := range moodValues {
		result = append(result, MoodDistributionItem{
			Mood:  mood,
			Count: counts[mood],
		})
	}

	return result
}

func createMoodNotes(logs []models.MoodLog) []MoodNote {
	result := []MoodNote{}

	for _, log := range logs {
		result = append(result, MoodNote{
			ID:        log.ID,
			Date:      formatDate(log.CreatedAt),
			Mood:      log.Mood,
			Note:      log.Note,
			Causes:    splitCauses(log.Causes),
			CreatedAt: log.CreatedAt.Format(time.RFC3339),
		})
	}

	return result
}

func calculateCauseSummaries(logs []models.MoodLog) []CauseSummary {
	moodValues := []int{1, 2, 3, 4, 5}
	causeMoodCounts := map[string]map[int]int{}

	for _, log := range logs {
		for _, cause := range splitCauses(log.Causes) {
			if causeMoodCounts[cause] == nil {
				causeMoodCounts[cause] = map[int]int{}
			}

			causeMoodCounts[cause][log.Mood]++
		}
	}

	result := []CauseSummary{}

	for cause, moodCounts := range causeMoodCounts {
		total := 0
		breakdown := []MoodDistributionItem{}

		for _, mood := range moodValues {
			count := moodCounts[mood]
			total += count

			breakdown = append(breakdown, MoodDistributionItem{
				Mood:  mood,
				Count: count,
			})
		}

		result = append(result, CauseSummary{
			Cause:         cause,
			TotalCount:    total,
			MoodBreakdown: breakdown,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalCount > result[j].TotalCount
	})

	return result
}

func splitCauses(causes string) []string {
	if causes == "" {
		return []string{}
	}

	rawCauses := strings.Split(causes, ",")
	result := []string{}

	for _, cause := range rawCauses {
		cause = strings.TrimSpace(cause)
		if cause != "" {
			result = append(result, cause)
		}
	}

	return result
}

func createMonthDateRange(month string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid month format")
	}

	end := start.AddDate(0, 1, -1)

	return start, end, nil
}

func createDateRange(startDate time.Time, endDate time.Time) []string {
	result := []string{}

	for current := startDate; !current.After(endDate); current = current.AddDate(0, 0, 1) {
		result = append(result, current.Format("2006-01-02"))
	}

	return result
}

func parseDate(value string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}

	return parsedDate, nil
}

func formatDate(value time.Time) string {
	return value.Format("2006-01-02")
}

func roundOneDecimal(value float64) float64 {
	result, _ := strconv.ParseFloat(strconv.FormatFloat(value, 'f', 1, 64), 64)
	return result
}
