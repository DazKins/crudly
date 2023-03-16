package service

import (
	"crudly/model"
	"database/sql"
	"fmt"
)

type postgresTableDeleter struct {
	postgres *sql.DB
}

func NewPostgresTableDeleter(postgres *sql.DB) postgresTableDeleter {
	return postgresTableDeleter{
		postgres,
	}
}

func (p postgresTableDeleter) DeleteTable(projectId model.ProjectId, name model.TableName) error {
	deleteSchemaQuery := getPostgresSchemaDeletionQuery(projectId, name)

	_, err := p.postgres.Exec(deleteSchemaQuery)

	if err != nil {
		return fmt.Errorf("error delete schema from postgres: %w", err)
	}

	deleteTableQuery := getPostgresDeleteTableQuery(projectId, name)

	_, err = p.postgres.Exec(deleteTableQuery)

	if err != nil {
		return fmt.Errorf("error delete table from postgres: %w", err)
	}

	return nil
}

func getPostgresDeleteTableQuery(projectId model.ProjectId, name model.TableName) string {
	return fmt.Sprintf(
		"DROP TABLE \"%s\"",
		getPostgresTableName(projectId, name),
	)
}

func getPostgresSchemaDeletionQuery(projectId model.ProjectId, name model.TableName) string {
	return fmt.Sprintf(
		"DELETE FROM \"%s\" WHERE name = '%s'",
		getPostgresSchemaTableName(projectId),
		name.String(),
	)
}
