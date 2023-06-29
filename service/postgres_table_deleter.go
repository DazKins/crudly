package service

import (
	"context"
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

func (p *postgresTableDeleter) DeleteTable(projectId model.ProjectId, name model.TableName) error {
	tx, err := p.postgres.BeginTx(context.Background(), nil)

	if err != nil {
		return fmt.Errorf("error opening postgres transaction: %w", err)
	}
	defer tx.Rollback()
	deleteSchemaQuery := getPostgresSchemaDeletionQuery(projectId, name)

	_, err = tx.Exec(deleteSchemaQuery)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	deleteTableQuery := getPostgresDeleteTableQuery(projectId, name)

	_, err = tx.Exec(deleteTableQuery)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error commiting postgres transaction: %w", err)
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
