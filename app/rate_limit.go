package app

import (
	"context"
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type callCountStore interface {
	IncrementCallCount(projectId model.ProjectId, ttl time.Duration) result.R[uint]
	GetCurrentCallCount(projectId model.ProjectId) result.R[uint]
}

type rateLimitStore interface {
	GetRateLimit(projectId model.ProjectId) result.R[uint]
}

type rateLimitManager struct {
	callCountStore       callCountStore
	rateLimitStore       rateLimitStore
	blockedProjectsCache map[model.ProjectId]struct{}
	mu                   sync.Mutex
}

func NewRateLimitManager(callCountStore callCountStore, rateLimitStore rateLimitStore) rateLimitManager {
	return rateLimitManager{
		callCountStore:       callCountStore,
		rateLimitStore:       rateLimitStore,
		blockedProjectsCache: map[model.ProjectId]struct{}{},
	}
}

const DEFAULT_RATE_LIMIT = uint(50_000)

func (r *rateLimitManager) GetDailyRateLimit(projectId model.ProjectId) result.R[uint] {
	rateLimitResult := r.rateLimitStore.GetRateLimit(projectId)

	if rateLimitResult.IsErr() {
		err := rateLimitResult.UnwrapErr()

		if _, ok := err.(errs.RateLimitNotFoundError); ok {
			return result.Ok(DEFAULT_RATE_LIMIT)
		}

		return result.Errf[uint]("error getting rate limit: %w", rateLimitResult.UnwrapErr())
	}

	return rateLimitResult
}

func (r *rateLimitManager) GetCurrentRateUsage(projectId model.ProjectId) result.R[uint] {
	return r.callCountStore.GetCurrentCallCount(projectId)
}

func (r *rateLimitManager) HandleUsage(projectId model.ProjectId) error {
	dailyRateLimit, currentRateusage := uint(0), uint(0)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		incrementRateUsageResult := r.callCountStore.IncrementCallCount(projectId, time.Hour*24)

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
