package service

import (
	"crudly/errs"
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

func (p *postgresEntityDeleter) DeleteEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) error {
	query := getPostgresDeleteEntityQuery(
		projectId,
		tableName,
		id,
	)

	res, err := p.postgres.Exec(query)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	count, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("error determining affected postgres rows: %w", err)
	}

	if count == 0 {
		return errs.EntityNotFoundError{}
	}

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
