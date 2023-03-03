package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util"
	"database/sql"
	"fmt"
	"strings"
)

type postgresTableFetcher struct {
	postgres *sql.DB
}

func NewPostgresTableFetcher(postgres *sql.DB) postgresTableFetcher {
	return postgresTableFetcher{
		postgres,
	}
}

func (p postgresTableFetcher) FetchTableSchema(projectId model.ProjectId, name model.TableName) util.Result[model.TableSchema] {
	query := "SELECT column_name, data_type " +
		"FROM information_schema.columns " +
		"WHERE table_name = '" + getPostgresTableName(projectId, name) + "'"

	rows, err := p.postgres.Query(query)

	if err != nil {
		return util.ResultErr[model.TableSchema](fmt.Errorf("error querying postgres: %w", err))
	}

	result := model.TableSchema{}

	defer rows.Close()

	rowCount := 0

	for rows.Next() {
		rowCount++

		columnName, dataType := "", ""
		err := rows.Scan(&columnName, &dataType)

		if err != nil {
			return util.ResultErr[model.TableSchema](fmt.Errorf("error scanning rows response: %w", err))
		}

		fieldDefinitionResult := getFieldTypeFromPostgresDataType(dataType)

		if fieldDefinitionResult.IsErr() {
			return util.ResultErr[model.TableSchema](fmt.Errorf("error getting field definition: %w", fieldDefinitionResult.UnwrapErr()))
		}

		result[columnName] = fieldDefinitionResult.Unwrap()
	}

	if rowCount == 0 {
		return util.ResultErr[model.TableSchema](errs.TableNotFoundError{})
	}

	return util.ResultOk(result)
}

func getFieldTypeFromPostgresDataType(dataType string) util.Result[model.FieldDefinition] {
	var resultType model.FieldType

	switch strings.ToLower(dataType) {
	case "uuid":
		resultType = model.FieldTypeId
	case "integer":
		resultType = model.FieldTypeInteger
	case "boolean":
		resultType = model.FieldTypeBoolean
	case "character varying":
		resultType = model.FieldTypeString
	case "timestamp without time zone":
		resultType = model.FieldTypeTime
	default:
		return util.ResultErr[model.FieldDefinition](fmt.Errorf("unsupported postgres datatype: %s", dataType))
	}

	return util.ResultOk(model.FieldDefinition{
		Type: resultType,
	})
}
