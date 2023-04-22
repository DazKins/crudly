package service

import (
	"crudly/model"
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
	id model.EntityId,
	partialEntity model.PartialEntity,
) error {
	query := getPostgresEntityUpdateQuery(
		projectId,
		tableName,
		id,
		partialEntity,
	)

	_, err := p.postgres.Query(query)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	return nil
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
		"UPDATE \"%s\" SET %s WHERE id = '%s'",
		getPostgresTableName(projectId, tableName),
		setQuery,
		id.String(),
	)
}
