package app

import (
	"crudly/model"
	"crudly/util/result"
	"time"
)

type rateLimiterStore interface {
	IncrementCallCount(projectId model.ProjectId, ttl time.Duration) error
	GetCurrentCallCount(projectId model.ProjectId) result.R[uint]
}

type rateLimitManager struct {
	rateLimiterStore rateLimiterStore
}

func NewRateLimitManager(rateLimiterStore rateLimiterStore) rateLimitManager {
	return rateLimitManager{
		rateLimiterStore,
	}
}

func (a rateLimitManager) GetDailyRateLimit(projectId model.ProjectId) uint {
	return 10 // TODO store and manage per project somewhere
}

func (a rateLimitManager) GetCurrentRateUsage(projectId model.ProjectId) result.R[uint] {
	return a.rateLimiterStore.GetCurrentCallCount(projectId)
}

func (a rateLimitManager) HandleUsage(projectId model.ProjectId) error {
	return a.rateLimiterStore.IncrementCallCount(projectId, time.Hour*24)
}
