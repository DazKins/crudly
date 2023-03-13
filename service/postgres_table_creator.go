package service

import (
	"crudly/model"
	"database/sql"
	"encoding/json"
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
	tableCreationQuery := getPostgresTableCreationQuery(
		projectId,
		name,
		schema,
	)

	res, err := p.postgres.Exec(tableCreationQuery)

	if err != nil {
		return fmt.Errorf("error creating postgres table: %w", err)
	}

	if res == nil {
		return fmt.Errorf("table %s already exists", name)
	}

	schemaCreationQuery := getPostgresTableSchemaCreationQuery(
		projectId,
		name,
		schema,
	)

	res, err = p.postgres.Exec(schemaCreationQuery)

	if err != nil {
		return fmt.Errorf("error creating postgres table: %w", err)
	}

	if res == nil {
		return fmt.Errorf("table %s already exists", name)
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
	for k, v := range schema {
		fieldQuery := getPostgresFieldQuery(k, v)
		query += fieldQuery + ","
	}
	query = strings.TrimSuffix(query, ",")
	query = query + ")"

	return query
}

func getPostgresFieldQuery(key string, fieldDefinition model.FieldDefinition) string {
	fieldQuery := key + " " + getPostgresDatatype(fieldDefinition.Type)

	if fieldDefinition.PrimaryKey {
		fieldQuery += " PRIMARY KEY"
	} else {
		fieldQuery += " NOT NULL"
	}

	return fieldQuery
}

func getPostgresDatatype(fieldType model.FieldType) string {
	switch fieldType {
	case model.FieldTypeId:
		return "uuid"
	case model.FieldTypeBoolean:
		return "boolean"
	case model.FieldTypeInteger:
		return "integer"
	case model.FieldTypeString:
		return "varchar"
	case model.FieldTypeTime:
		return "timestamp"
	case model.FieldTypeEnum:
		return "varchar"
	}
	panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldType))
}
