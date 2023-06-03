package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"database/sql"
	"fmt"
	"strings"
)

type postgresEntityUpdater struct {
	postgres *sql.DB
}

func NewPostgresEntityUpdater(postgres *sql.DB) postgresEntityUpdater {
	return postgresEntityUpdater{
		postgres,
	}
}

func (p postgresEntityUpdater) UpdateEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	id model.EntityId,
	partialEntity model.PartialEntity,
) result.R[model.Entity] {
	query := getPostgresEntityUpdateQuery(
		projectId,
		tableName,
		id,
		partialEntity,
	)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Errf[model.Entity]("error querying postgres: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return result.Err[model.Entity](errs.EntityNotFoundError{})
	}

	return parseEntityFromSqlRow(rows, tableSchema)
}

func getPostgresEntityUpdateQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
	partialEntity model.PartialEntity,
) string {
	setQuery := ""

	for k, v := range partialEntity {
		valResult := getPostgresFieldValue(v)

		if valResult.IsErr() {
			panic(fmt.Sprintf("error parsing field: %s: %s", k, valResult.UnwrapErr().Error()))
		}

		setQuery += fmt.Sprintf("\"%s\" = %s,", k, valResult.Unwrap())
	}

	setQuery = strings.TrimSuffix(setQuery, ",")

	return fmt.Sprintf(
		"UPDATE \"%s\" SET %s WHERE id = '%s' RETURNING *",
		getPostgresTableName(projectId, tableName),
		setQuery,
		id.String(),
	)
}
