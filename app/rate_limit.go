package app

import (
	"context"
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type rateLimiterStore interface {
	IncrementCallCount(projectId model.ProjectId, ttl time.Duration) result.R[uint]
	GetCurrentCallCount(projectId model.ProjectId) result.R[uint]
}

type rateLimitManager struct {
	rateLimiterStore     rateLimiterStore
	blockedProjectsCache map[model.ProjectId]struct{}
	mu                   sync.Mutex
}

func NewRateLimitManager(rateLimiterStore rateLimiterStore) rateLimitManager {
	return rateLimitManager{
		rateLimiterStore:     rateLimiterStore,
		blockedProjectsCache: map[model.ProjectId]struct{}{},
	}
}

func (r *rateLimitManager) GetDailyRateLimit(projectId model.ProjectId) result.R[uint] {
	return result.Ok(uint(10)) // TODO store and manage per project somewhere
}

func (r *rateLimitManager) GetCurrentRateUsage(projectId model.ProjectId) result.R[uint] {
	return r.rateLimiterStore.GetCurrentCallCount(projectId)
}

func (r *rateLimitManager) HandleUsage(projectId model.ProjectId) error {
	dailyRateLimit, currentRateusage := uint(0), uint(0)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		incrementRateUsageResult := r.rateLimiterStore.IncrementCallCount(projectId, time.Hour*24)

		if incrementRateUsageResult.IsErr() {
			return fmt.Errorf("error incrementing rate usage: %w", incrementRateUsageResult.UnwrapErr())
		}

		currentRateusage = incrementRateUsageResult.Unwrap()

		return nil
	})

	g.Go(func() error {
		dailyRateLimitResult := r.GetDailyRateLimit(projectId)

		if dailyRateLimitResult.IsErr() {
			return fmt.Errorf("error getting daily rate limit: %w", dailyRateLimitResult.UnwrapErr())
		}

		dailyRateLimit = dailyRateLimitResult.Unwrap()

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if currentRateusage >= dailyRateLimit {
		r.mu.Lock()
		r.blockedProjectsCache[projectId] = struct{}{}
		r.mu.Unlock()
	}

	return nil
}

func (r *rateLimitManager) ShouldBlockRequest(projectId model.ProjectId) bool {
	_, ok := r.blockedProjectsCache[projectId]

	return ok
}
