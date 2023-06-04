package dto

type RateLimitDto struct {
	DailyRateLimit   int `json:"dailyRateLimit"`
	CurrentRateUsage int `json:"currentRateUsage"`
}
