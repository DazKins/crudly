package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"database/sql"
	"encoding/json"
	"fmt"
)

type postgresTableFetcher struct {
	postgres *sql.DB
}

func NewPostgresTableFetcher(postgres *sql.DB) postgresTableFetcher {
	return postgresTableFetcher{
		postgres,
	}
}

func (p postgresTableFetcher) FetchTableSchema(
	projectId model.ProjectId,
	name model.TableName,
) result.Result[model.TableSchema] {
	query := "SELECT schema " +
		"FROM \"" + getPostgresSchemaTableName(projectId) + "\" " +
		"WHERE name = '" + name.String() + "'"

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.TableSchema](fmt.Errorf("error querying postgres: %w", err))
	}

	defer rows.Close()

	if !rows.Next() {
		return result.Err[model.TableSchema](errs.TableNotFoundError{})
	}

	schemaBytes := []byte{}

	rows.Scan(&schemaBytes)

	schema := model.TableSchema{}

	err = json.Unmarshal(schemaBytes, &schema)

	if err != nil {
		panic("error unmarshalling table schema")
	}

	return result.Ok(schema)
}
