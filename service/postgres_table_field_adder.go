package service

import (
	"context"
	"crudly/model"
	"crudly/util"
	"crudly/util/optional"
	"crudly/util/result"
	"database/sql"
	"fmt"
)

type postgresTableFieldAdder struct {
	postgres *sql.DB
}

func NewPostgresTableFieldAdder(postgres *sql.DB) postgresTableFieldAdder {
	return postgresTableFieldAdder{
		postgres,
	}
}

func (p *postgresTableFieldAdder) AddTableField(
	projectId model.ProjectId,
	tableName model.TableName,
	name model.FieldName,
	existingSchema model.TableSchema,
	definition model.FieldDefinition,
	defaultValue optional.O[any],
) error {
	tableQuery := ""

	if definition.IsOptional {
		tableQuery = getPostgresAddTableOptionalFieldQuery(
			projectId,
			tableName,
			name,
			definition,
		)
	} else {
		queryResult := getPostgresAddTableNonOptionalFieldQuery(
			projectId,
			tableName,
			name,
			definition,
			defaultValue.Unwrap(),
		)

		if queryResult.IsErr() {
			return queryResult.UnwrapErr()
		}

		tableQuery = queryResult.Unwrap()
	}

	newSchema := util.CopyMap(existingSchema)
	newSchema[name] = definition

	schemaUpdateQuery := getPostgresTableSchemaUpdateQuery(
		projectId,
		tableName,
		newSchema,
	)

	tx, err := p.postgres.BeginTx(context.Background(), nil)

	if err != nil {
		return fmt.Errorf("error opening postgres transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(tableQuery)

	if err != nil {
		return fmt.Errorf("unexpected error querying postgres: %w", err)
	}

	_, err = tx.Exec(schemaUpdateQuery)

	if err != nil {
		return fmt.Errorf("unexpected error querying postgres: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error commiting postgres transaction: %w", err)
	}

	return nil
}

func getPostgresAddTableOptionalFieldQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	name model.FieldName,
	definition model.FieldDefinition,
) string {
	return fmt.Sprintf(
		"ALTER TABLE \"%s\" ADD COLUMN \"%s\" %s",
		getPostgresTableName(projectId, tableName),
		name,
		getPostgresDatatype(definition.Type),
	)
}

func getPostgresAddTableNonOptionalFieldQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	name model.FieldName,
	definition model.FieldDefinition,
	defaultValue any,
) result.R[string] {
	postgresFieldValueResult := getPostgresFieldValue(defaultValue)

	if postgresFieldValueResult.IsErr() {
		return result.Errf[string]("couldnt get postgres field value for default value: %s", postgresFieldValueResult.UnwrapErr())
	}

	return result.Ok(
		fmt.Sprintf(
			"ALTER TABLE \"%s\" ADD COLUMN \"%s\" %s DEFAULT %s NOT NULL",
			getPostgresTableName(projectId, tableName),
			name,
			getPostgresDatatype(definition.Type),
			postgresFieldValueResult.Unwrap(),
		),
	)
}

func getPostgresTableSchemaUpdateQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	schema model.TableSchema,
) string {
	return fmt.Sprintf(
		"UPDATE \"%s\" SET schema = '%s' WHERE name = '%s'",
		getPostgresSchemaTableName(projectId),
		getSchemaJson(schema),
		tableName,
	)
}
