package service

import (
	"crudly/model"
)

func getPostgresTableName(projectId model.ProjectId, tableName model.TableName) string {
	return projectId.String() + "-table-" + tableName.String()
}

func getPostgresSchemaTableName(projectId model.ProjectId) string {
	return projectId.String() + "-tables"
}
