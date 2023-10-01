package dto

type RateLimitDto struct {
	DailyRateLimit   int `json:"dailyRateLimit"`
	CurrentRateUsage int `json:"currentRateUsage"`
}

type RateLimitUpdateRequest struct {
	DailyRateLimit uint `json:"dailyRateLimit"`
}
