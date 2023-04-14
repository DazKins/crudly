package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewAdminAuth() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(AdminContextKey) == nil {
				w.WriteHeader(404)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
