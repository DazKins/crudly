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

		fieldSchemaResult := getFieldSchemaFromPostgresDataType(dataType)

		if fieldSchemaResult.IsErr() {
			return util.ResultErr[model.TableSchema](fmt.Errorf("error getting field schema: %w", err))
		}

		result[columnName] = fieldSchemaResult.Unwrap()
	}

	if rowCount == 0 {
		return util.ResultErr[model.TableSchema](errs.TableNotFoundError{})
	}

	return util.ResultOk(result)
}

func getFieldSchemaFromPostgresDataType(dataType string) util.Result[model.FieldSchema] {
	switch strings.ToLower(dataType) {
	case "integer":
		return util.ResultOk(model.FieldSchemaInteger)
	case "boolean":
		return util.ResultOk(model.FieldSchemaBoolean)
	case "character varying":
		return util.ResultOk(model.FieldSchemaString)
	case "uuid":
		return util.ResultOk(model.FieldSchemaId)
	}
	return util.ResultErr[model.FieldSchema](fmt.Errorf("unsupported postgres datatype: %s", dataType))
}
