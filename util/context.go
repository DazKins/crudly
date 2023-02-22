package util

import (
	"crudly/model"
	"net/http"
)

type ContextKey string

const (
	ProjectIdContextKey = ContextKey("projectId")
)

func GetRequestProjectId(r *http.Request) model.ProjectId {
	return r.Context().Value(ProjectIdContextKey).(model.ProjectId)
}
