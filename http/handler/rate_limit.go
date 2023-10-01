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
	"io"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type rateLimitGetter interface {
	GetDailyRateLimit(projectId model.ProjectId) result.R[uint]
	GetCurrentRateUsage(projectId model.ProjectId) result.R[uint]
}

type rateLimitSetter interface {
	SetDailyRateLimit(projectId model.ProjectId, rateLimit uint) error
}

type rateLimitHandler struct {
	rateLimitGetter rateLimitGetter
	rateLimitSetter rateLimitSetter
}

func NewRateLimitHandler(rateLimitGetter rateLimitGetter, rateLimitSetter rateLimitSetter) rateLimitHandler {
	return rateLimitHandler{
		rateLimitGetter,
		rateLimitSetter,
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

func (rl *rateLimitHandler) PostRateLimit(w http.ResponseWriter, r *http.Request) {
	projectIdString := r.URL.Query().Get("projectId")

	if projectIdString == "" {
		w.WriteHeader(400)
		w.Write([]byte("projectId query param is required"))
		return
	}

	projectIdUuid, err := uuid.Parse(projectIdString)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("projectId query param must be a valid uuid"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var rateLimitUpdateRequestDto dto.RateLimitUpdateRequest
	err = json.Unmarshal(bodyBytes, &rateLimitUpdateRequestDto)

	if err != nil {
		middleware.AttachError(w, err)
		w.WriteHeader(400)
		w.Write([]byte("invalid request body"))
		return
	}

	err = rl.rateLimitSetter.SetDailyRateLimit(
		model.ProjectId(projectIdUuid),
		rateLimitUpdateRequestDto.DailyRateLimit,
	)

	if err != nil {
		middleware.AttachError(w, err)
		w.WriteHeader(500)
		w.Write([]byte("unexpected error setting rate limit"))
		return
	}
}
