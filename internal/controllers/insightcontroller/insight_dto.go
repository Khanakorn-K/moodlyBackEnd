package insightcontroller

type InsightResultResponse struct {
	TotalLogs        int                       `json:"totalLogs"`
	AverageMood      float64                   `json:"averageMood"`
	MoodDistribution map[string]int            `json:"moodDistribution"`
	CausesAnalysis   map[string]map[string]int `json:"causesAnalysis"`
}
