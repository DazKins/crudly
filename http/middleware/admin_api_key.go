package middleware

import (
	"context"
	"crudly/config"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAdminApiKey(config config.Config) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-api-key") == config.AdminApiKey {
				ctx := context.WithValue(r.Context(), AdminContextKey, struct{}{})
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
