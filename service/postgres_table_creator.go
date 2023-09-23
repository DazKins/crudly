package service

import (
	"context"
	"crudly/model"
	"database/sql"
	"encoding/json"
	"fmt"
)

type postgresTableCreator struct {
	postgres *sql.DB
}

func NewPostgresTableCreator(postgres *sql.DB) postgresTableCreator {
	return postgresTableCreator{
		postgres,
	}
}

func (p *postgresTableCreator) CreateTable(
	projectId model.ProjectId,
	name model.TableName,
	schema model.TableSchema,
) error {
	tx, err := p.postgres.BeginTx(context.Background(), nil)

	if err != nil {
		return fmt.Errorf("error opening postgres transaction: %w", err)
	}
	defer tx.Rollback()

	tableCreationQuery := getPostgresTableCreationQuery(
		projectId,
		name,
		schema,
	)

	_, err = tx.Exec(tableCreationQuery)

	if err != nil {
		return fmt.Errorf("error creating postgres table: %w", err)
	}

	schemaCreationQuery := getPostgresTableSchemaCreationQuery(
		projectId,
		name,
		schema,
	)

	_, err = tx.Exec(schemaCreationQuery)

	if err != nil {
		return fmt.Errorf("error creating postgres table: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error commiting postgres transaction: %w", err)
	}

	return nil
}

func getPostgresTableSchemaCreationQuery(
	projectId model.ProjectId,
	name model.TableName,
	schema model.TableSchema,
) string {
	return fmt.Sprintf(
		"INSERT INTO \"%s\"(%s, %s) VALUES ('%s', '%s')",
		getPostgresSchemaTableName(projectId),
		"name", "schema",
		name, getSchemaJson(schema),
	)
}

func getSchemaJson(schema model.TableSchema) string {
	json, err := json.Marshal(schema)

	if err != nil {
		panic("error marshalling table schema json")
	}

	return string(json)
}

func getPostgresTableCreationQuery(
	projectId model.ProjectId,
	name model.TableName,
	schema model.TableSchema,
) string {
	query := "CREATE TABLE \"" + getPostgresTableName(projectId, name) + "\"("
	query += getPostgresFieldQuery("id", model.FieldDefinition{
		Type:       model.FieldTypeId,
		PrimaryKey: true,
	})
	for k, v := range schema {
		query += "," + getPostgresFieldQuery(k, v)
	}
	query = query + ")"

	return query
}

func getPostgresFieldQuery(key model.FieldName, fieldDefinition model.FieldDefinition) string {
	fieldQuery := fmt.Sprintf(
		"\"%s\" %s",
		key.String(),
		getPostgresDatatype(fieldDefinition.Type),
	)

	if fieldDefinition.PrimaryKey {
		fieldQuery += " PRIMARY KEY"
	} else if !fieldDefinition.IsOptional {
		fieldQuery += " NOT NULL"
	}

	return fieldQuery
}
