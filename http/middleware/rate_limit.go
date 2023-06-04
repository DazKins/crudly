package middleware

import (
	"crudly/ctx"
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type RateLimitGetter interface {
	GetDailyRateLimit(projectId model.ProjectId) uint
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
				currentUsageResult := rateLimitGetter.GetCurrentRateUsage(projectId)

				if currentUsageResult.IsErr() {
					fmt.Printf("error getting rate limit usage: %s\n", currentUsageResult.UnwrapErr().Error())
				} else {
					currentUsage := currentUsageResult.Unwrap()

					if currentUsage >= rateLimitGetter.GetDailyRateLimit(projectId) {
						blockedProjects[projectId] = struct{}{}
					}
				}

				err := rateLimitHandler.HandleUsage(projectId)

				if err != nil {
					fmt.Printf("error handle rate limit usage: %s\n", err)
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}
