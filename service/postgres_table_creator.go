package service

import (
	"crudly/model"
	"database/sql"
	"fmt"
	"strings"
)

type postgresTableCreator struct {
	postgres *sql.DB
}

func NewPostgresTableCreator(postgres *sql.DB) postgresTableCreator {
	return postgresTableCreator{
		postgres,
	}
}

func (p postgresTableCreator) CreateTable(
	projectId model.ProjectId,
	name model.TableName,
	schema model.TableSchema,
) error {
	query := getPostgresTableCreationQuery(
		projectId,
		name,
		schema,
	)

	res, err := p.postgres.Exec(query)

	if res == nil {
		return fmt.Errorf("table %s already exists", name)
	}

	if err != nil {
		return fmt.Errorf("error creating postgres table: %w", err)
	}

	return nil
}

func getPostgresTableCreationQuery(
	projectId model.ProjectId,
	name model.TableName,
	schema model.TableSchema,
) string {
	query := "CREATE TABLE \"" + getPostgresTableName(projectId, name) + "\"("
	for k, v := range schema {
		fieldQuery := getPostgresFieldQuery(k, v)
		query += fieldQuery + ","
	}
	query = strings.TrimSuffix(query, ",")
	query = query + ")"

	return query
}

func getPostgresFieldQuery(key string, schema model.FieldSchema) string {
	return key + " " + getPostgresDatatype(schema) + " NOT NULL"
}

func getPostgresDatatype(schema model.FieldSchema) string {
	switch schema {
	case model.FieldSchemaBoolean:
		return "boolean"
	case model.FieldSchemaInteger:
		return "integer"
	case model.FieldSchemaString:
		return "varchar"
	case model.FieldSchemaId:
		return "uuid"
	}
	panic(fmt.Sprintf("invalid field schema has entered the system: %+v", schema))
}
