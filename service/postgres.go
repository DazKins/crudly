package service

import "crudly/model"

func getPostgresTableName(projectId model.ProjectId, tableName model.TableName) string {
	return projectId.String() + "-" + tableName.String()
}
