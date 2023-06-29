package handler

import (
	"crudly/ctx"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util"
	"crudly/util/result"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type rateLimitGetter interface {
	GetDailyRateLimit(projectId model.ProjectId) result.R[uint]
	GetCurrentRateUsage(projectId model.ProjectId) result.R[uint]
}

type rateLimitHandler struct {
	rateLimitGetter rateLimitGetter
}

func NewRateLimitHandler(rateLimitGetter rateLimitGetter) rateLimitHandler {
	return rateLimitHandler{
		rateLimitGetter,
	}
}

func (rl *rateLimitHandler) GetRateLimit(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	dailyRateLimit, currentRateusage := uint(0), uint(0)

	g, _ := errgroup.WithContext(r.Context())

	g.Go(func() error {
		currentRateUsageResult := rl.rateLimitGetter.GetCurrentRateUsage(projectId)

		if currentRateUsageResult.IsErr() {
			return fmt.Errorf("error getting current rate usage: %w", currentRateUsageResult.UnwrapErr())
		}

		currentRateusage = currentRateUsageResult.Unwrap()

		return nil
	})

	g.Go(func() error {
		dailyRateLimitResult := rl.rateLimitGetter.GetDailyRateLimit(projectId)

		if dailyRateLimitResult.IsErr() {
			return fmt.Errorf("error getting daily rate limit: %w", dailyRateLimitResult.UnwrapErr())
		}

		dailyRateLimit = dailyRateLimitResult.Unwrap()

		return nil
	})

	if err := g.Wait(); err != nil {
		middleware.AttachError(w, err)
		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting rate limit"))
		return
	}

	resBodyBytes, _ := json.Marshal(dto.RateLimitDto{
		DailyRateLimit:   int(dailyRateLimit),
		CurrentRateUsage: int(util.Min(dailyRateLimit, currentRateusage)),
	})

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}
