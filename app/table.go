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
	) result.Result[model.TableSchema]
}

type tableSchemaValidator interface {
	ValidateTableSchema(schema model.TableSchema) error
}

type tableManager struct {
	tableSchemaFetcher   tableSchemaFetcher
	tableCreator         tableCreator
	tableSchemaValidator tableSchemaValidator
}

func NewTableManager(
	tableSchemaFetcher tableSchemaFetcher,
	tableCreator tableCreator,
	tableSchemaValidator tableSchemaValidator,
) tableManager {
	return tableManager{
		tableSchemaFetcher,
		tableCreator,
		tableSchemaValidator,
	}
}

func (t tableManager) GetTableSchema(projectId model.ProjectId, name model.TableName) result.Result[model.TableSchema] {
	tableSchemaResult := t.tableSchemaFetcher.FetchTableSchema(projectId, name)

	if tableSchemaResult.IsErr() {
		err := tableSchemaResult.UnwrapErr()

		if _, ok := err.(errs.TableNotFoundError); ok {
			return tableSchemaResult
		}

		return result.Err[model.TableSchema](fmt.Errorf("error fetching table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	tableSchema := tableSchemaResult.Unwrap()

	delete(tableSchema, "id")

	return result.Ok(tableSchema)
}

func (t tableManager) CreateTable(projectId model.ProjectId, name model.TableName, schema model.TableSchema) error {
	_, ok := schema["id"]

	if ok {
		return errs.IdFieldAlreadyExistsError{}
	}

	schema["id"] = model.FieldDefinition{
		Type:       model.FieldTypeId,
		PrimaryKey: true,
	}

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
