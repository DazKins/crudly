package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type postgresEntityFetcher struct {
	postgres *sql.DB
}

func NewPostgresEntityFetcher(postgres *sql.DB) postgresEntityFetcher {
	return postgresEntityFetcher{
		postgres,
	}
}

func (p postgresEntityFetcher) FetchEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	id model.EntityId,
) result.Result[model.Entity] {
	query := getPostgresEntityQuery(
		projectId,
		tableName,
		id,
	)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.Entity](fmt.Errorf("error querying postgres: %w", err))
	}

	defer rows.Close()

	if !rows.Next() {
		return result.Err[model.Entity](errs.EntityNotFoundError{})
	}

	entity := model.Entity{}

	columnTypes, _ := rows.ColumnTypes()
	columns := make([]any, len(columnTypes))

	for i := range columns {
		columns[i] = new(string)
	}

	err = rows.Scan(columns...)

	for i, column := range columns {
		str := *(column.(*string))
		columnType := *columnTypes[i]

		fieldDefinition := tableSchema[columnType.Name()]

		entity[columnType.Name()] = parsePostgresFieldString(str, fieldDefinition.Type)
	}

	if err != nil {
		return result.Err[model.Entity](fmt.Errorf("error scaning postgres rows: %w", err))
	}

	return result.Ok(entity)
}

func (p postgresEntityFetcher) FetchEntities(
	projectId model.ProjectId,
	tableName model.TableName,
	tableSchema model.TableSchema,
	paginationParams model.PaginationParams,
) result.Result[model.Entities] {
	query := getPostgresEntitiesQuery(
		projectId,
		tableName,
		paginationParams,
	)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.Entities](fmt.Errorf("error querying postgres: %w", err))
	}

	defer rows.Close()

	columnTypes, _ := rows.ColumnTypes()
	columns := make([]any, len(columnTypes))

	entities := model.Entities{}

	for i := range columns {
		columns[i] = new(string)
	}

	for rows.Next() {
		entity := model.Entity{}

		err = rows.Scan(columns...)

		for i, column := range columns {
			str := *(column.(*string))
			columnType := *columnTypes[i]

			fieldDefinition := tableSchema[columnType.Name()]

			entity[columnType.Name()] = parsePostgresFieldString(str, fieldDefinition.Type)
		}

		if err != nil {
			return result.Err[model.Entities](fmt.Errorf("error scaning postgres rows: %w", err))
		}

		entities = append(entities, entity)
	}

	return result.Ok(entities)
}

func parsePostgresFieldString(str string, fieldType model.FieldType) any {
	switch fieldType {
	case model.FieldTypeId:
		uuid, err := uuid.Parse(str)

		if err != nil {
			panic("couldn't parse uuid in sql response")
		}

		return uuid
	case model.FieldTypeInteger:
		integer, err := strconv.Atoi(str)

		if err != nil {
			panic("couldn't parse integer in sql response")
		}

		return integer
	case model.FieldTypeBoolean:
		if strings.ToLower(str) == "false" {
			return false
		}
		return true
	case model.FieldTypeString:
		return str
	case model.FieldTypeTime:
		time, err := time.Parse(PostgresTimeFormat, str)

		if err != nil {
			panic(fmt.Sprintf("unexpected time format: %s", str))
		}

		return time
	}

	return nil
}

func getPostgresEntityQuery(projectId model.ProjectId, tableName model.TableName, id model.EntityId) string {
	return "SELECT * FROM \"" + getPostgresTableName(projectId, tableName) + "\" WHERE id = '" + id.String() + "'"
}

func getPostgresEntitiesQuery(projectId model.ProjectId, tableName model.TableName, paginationParams model.PaginationParams) string {
	return "SELECT * FROM \"" + getPostgresTableName(projectId, tableName) + "\"" +
		" LIMIT " + paginationParams.Limit.String() +
		" OFFSET " + paginationParams.Offset.String()
}
