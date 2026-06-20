package services

import (
	"errors"
	"moodly/internal/domain/entities"
	"moodly/internal/domain/repositories"
	"moodly/utils"
	"sort"
	"strings"
	"time"
)

type OverviewService struct {
	repo repositories.OverviewRepositoryInterface
}

func NewOverviewService(repo repositories.OverviewRepositoryInterface) *OverviewService {
	return &OverviewService{repo: repo}
}

type DailyCauseDistribution struct {
	Cause string
	Count int
}

type DailyMoodAverage struct {
	Date              string
	AverageMood       float64
	TotalLogs         int
	CauseDistribution []DailyCauseDistribution
}

type MoodDistribution struct {
	Mood  int
	Count int
}

type MoodNote struct {
	ID        uint
	Date      string
	Mood      int
	Note      string
	Causes    []string
	CreatedAt string
}

type CauseSummary struct {
	Cause         string
	TotalCount    int
	MoodBreakdown []MoodDistribution
}

type MonthlyAverageMoodResult struct {
	Month             string
	StartDate         string
	EndDate           string
	TotalLogs         int
	AverageMood       float64
	DailyMoodAverages []DailyMoodAverage
}

type OverviewResult struct {
	StartDate         string
	EndDate           string
	TotalLogs         int
	AverageMood       float64
	DailyMoodAverages []DailyMoodAverage
	MoodDistribution  []MoodDistribution
	MoodNotes         []MoodNote
	CauseSummaries    []CauseSummary
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
		AverageMood:       utils.CalculateAverageMood(logs),
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

	parsedStart, err := parseDate(startDate)
	if err != nil {
		return nil, err
	}

	parsedEnd, err := parseDate(endDate)
	if err != nil {
		return nil, err
	}

	if parsedStart.After(parsedEnd) {
		return nil, errors.New("invalid date range")
	}

	logs, err := s.repo.FindMoodLogsByDateRange(userID, parsedStart, parsedEnd)
	if err != nil {
		return nil, err
	}

	dateRange := createDateRange(parsedStart, parsedEnd)

	return &OverviewResult{
		StartDate:         formatDate(parsedStart),
		EndDate:           formatDate(parsedEnd),
		TotalLogs:         len(logs),
		AverageMood:       utils.CalculateAverageMood(logs),
		DailyMoodAverages: calculateDailyMoodAverages(logs, dateRange),
		MoodDistribution:  calculateMoodDistribution(logs),
		MoodNotes:         createMoodNotes(logs),
		CauseSummaries:    calculateCauseSummaries(logs),
	}, nil
}

func calculateDailyMoodAverages(logs []entities.MoodLogEntity, dateRange []string) []DailyMoodAverage {
	type dailyAccum struct {
		TotalMood   int
		TotalLogs   int
		CauseCounts map[string]int
	}

	dailyMap := map[string]*dailyAccum{}

	for _, log := range logs {
		date := formatDate(log.CreatedAt)
		if dailyMap[date] == nil {
			dailyMap[date] = &dailyAccum{CauseCounts: map[string]int{}}
		}
		dailyMap[date].TotalMood += log.Mood
		dailyMap[date].TotalLogs++
		for _, cause := range splitCauses(log.Causes) {
			dailyMap[date].CauseCounts[cause]++
		}
	}

	result := make([]DailyMoodAverage, 0, len(dateRange))

	for _, date := range dateRange {
		accum := dailyMap[date]
		if accum == nil {
			result = append(result, DailyMoodAverage{
				Date:              date,
				CauseDistribution: []DailyCauseDistribution{},
			})
			continue
		}

		causes := make([]DailyCauseDistribution, 0, len(accum.CauseCounts))
		for cause, count := range accum.CauseCounts {
			causes = append(causes, DailyCauseDistribution{Cause: cause, Count: count})
		}
		sort.Slice(causes, func(i, j int) bool {
			return causes[i].Count > causes[j].Count
		})

		result = append(result, DailyMoodAverage{
			Date:              date,
			AverageMood:       utils.RoundOneDecimal(float64(accum.TotalMood) / float64(accum.TotalLogs)),
			TotalLogs:         accum.TotalLogs,
			CauseDistribution: causes,
		})
	}

	return result
}

func calculateMoodDistribution(logs []entities.MoodLogEntity) []MoodDistribution {
	counts := map[int]int{}
	for _, log := range logs {
		counts[log.Mood]++
	}

	result := make([]MoodDistribution, 0, 5)
	for _, mood := range []int{1, 2, 3, 4, 5} {
		result = append(result, MoodDistribution{Mood: mood, Count: counts[mood]})
	}
	return result
}

func createMoodNotes(logs []entities.MoodLogEntity) []MoodNote {
	result := make([]MoodNote, 0, len(logs))
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

func calculateCauseSummaries(logs []entities.MoodLogEntity) []CauseSummary {
	causeMoodCounts := map[string]map[int]int{}

	for _, log := range logs {
		for _, cause := range splitCauses(log.Causes) {
			if causeMoodCounts[cause] == nil {
				causeMoodCounts[cause] = map[int]int{}
			}
			causeMoodCounts[cause][log.Mood]++
		}
	}

	result := make([]CauseSummary, 0, len(causeMoodCounts))
	for cause, moodCounts := range causeMoodCounts {
		total := 0
		breakdown := make([]MoodDistribution, 0, 5)
		for _, mood := range []int{1, 2, 3, 4, 5} {
			count := moodCounts[mood]
			total += count
			breakdown = append(breakdown, MoodDistribution{Mood: mood, Count: count})
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
	raw := strings.Split(causes, ",")
	result := make([]string, 0, len(raw))
	for _, c := range raw {
		if c = strings.TrimSpace(c); c != "" {
			result = append(result, c)
		}
	}
	return result
}

func createMonthDateRange(month string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("invalid month format")
	}
	return start, start.AddDate(0, 1, -1), nil
}

func createDateRange(start, end time.Time) []string {
	result := []string{}
	for cur := start; !cur.After(end); cur = cur.AddDate(0, 0, 1) {
		result = append(result, cur.Format("2006-01-02"))
	}
	return result
}

func parseDate(value string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New("invalid date format")
	}
	return t, nil
}

func formatDate(value time.Time) string {
	return value.Format("2006-01-02")
}
