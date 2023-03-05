package ctx

import (
	"crudly/model"
	"net/http"
)

type ContextKey string

const (
	ProjectIdContextKey = ContextKey("projectId")
	TableNameContextKey = ContextKey("tableName")
)

func GetRequestProjectId(r *http.Request) model.ProjectId {
	return r.Context().Value(ProjectIdContextKey).(model.ProjectId)
}

func GetRequestTableName(r *http.Request) model.TableName {
	return r.Context().Value(TableNameContextKey).(model.TableName)
}
