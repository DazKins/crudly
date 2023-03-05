package service

import (
	"crudly/model"
	"database/sql"
	"fmt"
)

type postgresEntityDeleter struct {
	postgres *sql.DB
}

func NewPostgresEntityDeleter(postgres *sql.DB) postgresEntityDeleter {
	return postgresEntityDeleter{
		postgres,
	}
}

func (p postgresEntityDeleter) DeleteEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) error {
	query := getPostgresDeleteEntityQuery(
		projectId,
		tableName,
		id,
	)

	res, err := p.postgres.Query(query)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	test := res.Next()

	fmt.Printf("%v", test)

	return nil
}

func getPostgresDeleteEntityQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) string {
	return fmt.Sprintf(
		"DELETE FROM \"%s\" WHERE id = '%s'",
		getPostgresTableName(projectId, tableName),
		id.String(),
	)
}
