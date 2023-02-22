package app

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util"
)

type tableCreator interface {
	CreateTable(
		projectId model.ProjectId,
		name model.TableName,
		schema model.TableSchema,
	) error
}

type tableSchemaFetcher interface {
	FetchTableSchema(
		projectId model.ProjectId,
		name model.TableName,
	) util.Result[model.TableSchema]
}

type tableManager struct {
	tableSchemaFetcher tableSchemaFetcher
	tableCreator       tableCreator
}

func NewTableManager(tableSchemaFetcher tableSchemaFetcher, tableCreator tableCreator) tableManager {
	return tableManager{
		tableSchemaFetcher,
		tableCreator,
	}
}

func (t tableManager) GetTableSchema(projectId model.ProjectId, name model.TableName) util.Result[model.TableSchema] {
	return t.tableSchemaFetcher.FetchTableSchema(projectId, name)
}

func (t tableManager) CreateTable(projectId model.ProjectId, name model.TableName, schema model.TableSchema) error {
	_, ok := schema["id"]

	if ok {
		return errs.IdFieldAlreadyExistsError{}
	}

	schema["id"] = model.FieldSchemaId

	return t.tableCreator.CreateTable(
		projectId,
		name,
		schema,
	)
}
