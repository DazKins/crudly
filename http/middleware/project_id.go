package middleware

import (
	"context"
	"crudly/ctx"
	"crudly/http/dto"
	"net/http"
)

type projectId struct{}

func NewProjectId() projectId {
	return projectId{}
}

func (projectId) Attach(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		projectIdDto := dto.ProjectIdDto(r.Header.Get("x-project-id"))
		projectIdResult := projectIdDto.ToModel()

		if projectIdResult.IsErr() {
			AttachError(w, projectIdResult.UnwrapErr())
			w.WriteHeader(400)
			w.Write([]byte("invalid project id header"))
			return
		}

		ctx := context.WithValue(r.Context(), ctx.ProjectIdContextKey, projectIdResult.Unwrap())

		h(w, r.WithContext(ctx))
	}
}
