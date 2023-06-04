package middleware

import (
	"crudly/ctx"
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

type RateLimitGetter interface {
	GetDailyRateLimit(projectId model.ProjectId) result.R[uint]
	GetCurrentRateUsage(projectId model.ProjectId) result.R[uint]
}

type RateLimitHandler interface {
	HandleUsage(projectId model.ProjectId) error
}

func NewRateLimit(rateLimitGetter RateLimitGetter, rateLimitHandler RateLimitHandler) mux.MiddlewareFunc {
	blockedProjects := map[model.ProjectId]struct{}{}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(AdminContextKey) != nil {
				h.ServeHTTP(w, r)
				return
			}

			projectId := ctx.GetRequestProjectId(r)

			if _, ok := blockedProjects[projectId]; ok {
				w.WriteHeader(429)
				w.Write([]byte("rate limit exceeded"))
				return
			}

			go func() {
				dailyRateLimit, currentRateusage := uint(0), uint(0)

				g, _ := errgroup.WithContext(r.Context())

				g.Go(func() error {
					currentRateUsageResult := rateLimitGetter.GetCurrentRateUsage(projectId)

					if currentRateUsageResult.IsErr() {
						return fmt.Errorf("error getting current rate usage: %w", currentRateUsageResult.UnwrapErr())
					}

					currentRateusage = currentRateUsageResult.Unwrap()

					return nil
				})

				g.Go(func() error {
					dailyRateLimitResult := rateLimitGetter.GetDailyRateLimit(projectId)

					if dailyRateLimitResult.IsErr() {
						return fmt.Errorf("error getting daily rate limit: %w", dailyRateLimitResult.UnwrapErr())
					}

					dailyRateLimit = dailyRateLimitResult.Unwrap()

					return nil
				})

				if err := g.Wait(); err != nil {
					fmt.Printf("%s\n", err.Error())
					return
				}

				if currentRateusage >= dailyRateLimit {
					blockedProjects[projectId] = struct{}{}
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}
