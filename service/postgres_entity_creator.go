package service

import (
	"crudly/model"
	"crudly/util"
	"database/sql"
	"fmt"
)

type postgresEntityCreator struct {
	postgres *sql.DB
}

func NewPostgresEntityCreator(postgres *sql.DB) postgresEntityCreator {
	return postgresEntityCreator{
		postgres,
	}
}

func (p postgresEntityCreator) CreateEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	id model.EntityId,
	entity model.Entity,
) error {
	query := getPostgresCreateEntityQuery(
		projectId,
		tableName,
		tableSchema,
		id,
		entity,
	)

	_, err := p.postgres.Query(query)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	return nil
}

func getPostgresCreateEntityQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	id model.EntityId,
	entity model.Entity,
) string {
	query := "INSERT INTO \"" + getPostgresTableName(projectId, tableName) + "\"("

	keys := util.GetMapKeys(entity)

	for _, k := range keys {
		query = query + k + ","
	}

	query += "id) VALUES ("

	for _, k := range keys {
		query += fmt.Sprintf("'%+v'", entity[k]) + ","
	}

	query += "'" + id.String() + "')"
	return query
}
