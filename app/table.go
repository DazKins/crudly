package app

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"fmt"
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
	) result.R[model.TableSchema]

	FetchTableSchemas(
		projectId model.ProjectId,
	) result.R[model.TableSchemas]
}

type tableDeleter interface {
	DeleteTable(
		projectId model.ProjectId,
		name model.TableName,
	) error
}

type tableSchemaValidator interface {
	ValidateTableSchema(schema model.TableSchema) error
}

type entityCountFetcher interface {
	FetchTotalEntityCount(
		projectId model.ProjectId,
		tableName model.TableName,
	) result.R[uint]
}

type tableManager struct {
	tableSchemaFetcher   tableSchemaFetcher
	tableCreator         tableCreator
	tableDeleter         tableDeleter
	tableSchemaValidator tableSchemaValidator
	entityCountFetcher   entityCountFetcher
}

func NewTableManager(
	tableSchemaFetcher tableSchemaFetcher,
	tableCreator tableCreator,
	tableDeleter tableDeleter,
	tableSchemaValidator tableSchemaValidator,
	entityCountFetcher entityCountFetcher,
) tableManager {
	return tableManager{
		tableSchemaFetcher,
		tableCreator,
		tableDeleter,
		tableSchemaValidator,
		entityCountFetcher,
	}
}

func (t *tableManager) GetTableSchema(projectId model.ProjectId, name model.TableName) result.R[model.TableSchema] {
	tableSchemaResult := t.tableSchemaFetcher.FetchTableSchema(projectId, name)

	if tableSchemaResult.IsErr() {
		err := tableSchemaResult.UnwrapErr()

		if _, ok := err.(errs.TableNotFoundError); ok {
			return tableSchemaResult
		}

		return result.Err[model.TableSchema](fmt.Errorf("error fetching table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	tableSchema := tableSchemaResult.Unwrap()

	return result.Ok(tableSchema)
}

func (t *tableManager) GetTableSchemas(projectId model.ProjectId) result.R[model.TableSchemas] {
	tableSchemasResult := t.tableSchemaFetcher.FetchTableSchemas(projectId)

	if tableSchemasResult.IsErr() {
		return result.Errf[model.TableSchemas]("error fetching table schemas: %w", tableSchemasResult.UnwrapErr())
	}

	return result.Ok(tableSchemasResult.Unwrap())
}

func (t *tableManager) CreateTable(projectId model.ProjectId, name model.TableName, schema model.TableSchema) error {
	err := t.tableSchemaValidator.ValidateTableSchema(schema)

	if err != nil {
		return errs.NewInvalidTableError(err)
	}

	return t.tableCreator.CreateTable(
		projectId,
		name,
		schema,
	)
}

func (t *tableManager) DeleteTable(projectId model.ProjectId, name model.TableName) error {
	return t.tableDeleter.DeleteTable(projectId, name)
}

func (t *tableManager) GetTotalEntityCount(
	projectId model.ProjectId,
	tableName model.TableName,
) result.R[uint] {
	return t.entityCountFetcher.FetchTotalEntityCount(projectId, tableName)
}
