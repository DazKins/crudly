package service

import (
	"context"
	"crudly/model"
	"crudly/util"
	"database/sql"
	"fmt"
)

type postgresTableFieldDeleter struct {
	postgres *sql.DB
}

func NewPostgresTableFieldDeleter(postgres *sql.DB) postgresTableFieldDeleter {
	return postgresTableFieldDeleter{
		postgres,
	}
}

func (p *postgresTableFieldDeleter) DeleteField(
	projectId model.ProjectId,
	tableName model.TableName,
	existingSchema model.TableSchema,
	fieldName model.FieldName,
) error {
	deleteFieldQuery := getPostgresTableFieldDeleteQuery(
		projectId,
		tableName,
		fieldName,
	)

	tx, err := p.postgres.BeginTx(context.Background(), nil)

	if err != nil {
		return fmt.Errorf("error opening postgres transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(deleteFieldQuery)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	newSchema := util.CopyMap(existingSchema)
	delete(newSchema, fieldName)

	schemaUpdateQuery := getPostgresTableSchemaUpdateQuery(
		projectId,
		tableName,
		newSchema,
	)

	_, err = tx.Exec(schemaUpdateQuery)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error commiting postgres transaction: %w", err)
	}

	return nil
}

func getPostgresTableFieldDeleteQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	fieldName model.FieldName,
) string {
	return fmt.Sprintf(
		"ALTER TABLE \"%s\" DROP COLUMN \"%s\"",
		getPostgresTableName(projectId, tableName),
		fieldName,
	)
}
