package middleware

import (
	"crudly/ctx"
	"crudly/model"
	"net/http"

	"github.com/gorilla/mux"
)

type RateLimitHandler interface {
	HandleUsage(projectId model.ProjectId) error
	ShouldBlockRequest(projectId model.ProjectId) bool
}

func NewRateLimit(rateLimitHandler RateLimitHandler) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(AdminContextKey) != nil {
				h.ServeHTTP(w, r)
				return
			}

			projectId := ctx.GetRequestProjectId(r)

			if rateLimitHandler.ShouldBlockRequest(projectId) {
				w.WriteHeader(429)
				w.Write([]byte("rate limit exceeded"))
				return
			}

			go func() {
				rateLimitHandler.HandleUsage(projectId)
			}()

			h.ServeHTTP(w, r)
		})
	}
}
