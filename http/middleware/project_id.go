package middleware

import (
	"context"
	"crudly/http/dto"
	"crudly/util"
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
			w.WriteHeader(400)
			return
		}

		ctx := context.WithValue(r.Context(), util.ProjectIdContextKey, projectIdResult.Unwrap())

		h(w, r.WithContext(ctx))
	}
}
