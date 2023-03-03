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
		fieldQuery := getPostgresFieldQuery(k, v.Type)
		query += fieldQuery + ","
	}
	query = strings.TrimSuffix(query, ",")
	query = query + ")"

	return query
}

func getPostgresFieldQuery(key string, fieldType model.FieldType) string {
	return key + " " + getPostgresDatatype(fieldType) + " NOT NULL"
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
	}
	panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldType))
}
