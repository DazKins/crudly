package middleware

import (
	"context"
	"crudly/ctx"
	"crudly/http/dto"
	"net/http"

	"github.com/gorilla/mux"
)

func NewProjectId() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			projectIdDto := dto.ProjectIdDto(r.Header.Get("x-project-id"))
			projectIdResult := projectIdDto.ToModel()

			if projectIdResult.IsErr() {
				AttachError(w, projectIdResult.UnwrapErr())
				w.WriteHeader(400)
				w.Write([]byte("invalid project id header"))
				return
			}

			ctx := context.WithValue(r.Context(), ctx.ProjectIdContextKey, projectIdResult.Unwrap())

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
