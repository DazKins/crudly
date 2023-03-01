package service

import (
	"crudly/model"
	"crudly/util"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type postgresEntityCreator struct {
	postgres *sql.DB
}

func NewPostgresEntityCreator(postgres *sql.DB) postgresEntityCreator {
	return postgresEntityCreator{
		postgres,
	}
}

func (p postgresEntityCreator) CreateEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
	entity model.Entity,
) error {
	query := getPostgresCreateEntityQuery(
		projectId,
		tableName,
		id,
		entity,
	)

	_, err := p.postgres.Query(query)

	if err != nil {
		return fmt.Errorf("error querying postgres: %w", err)
	}

	return nil
}

func getPostgresCreateEntityQuery(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
	entity model.Entity,
) string {
	query := "INSERT INTO \"" + getPostgresTableName(projectId, tableName) + "\"("

	keys := util.GetMapKeys(entity)

	for _, k := range keys {
		query = query + k + ","
	}

	query += "id) VALUES ("

	for _, k := range keys {
		postgresFieldValueResult := getPostgresFieldValue(entity[k])

		if postgresFieldValueResult.IsErr() {
			panic(fmt.Sprintf("error parsing field: %s: %s", k, postgresFieldValueResult.UnwrapErr().Error()))
		}

		query += postgresFieldValueResult.Unwrap() + ","
	}

	query += "'" + id.String() + "')"
	return query
}

const PostgresTimeFormat = "2006-01-02T15:04:05Z"

func getPostgresFieldValue(field any) util.Result[string] {
	switch field.(type) {
	case uuid.UUID:
		return util.ResultOk("'" + field.(uuid.UUID).String() + "'")
	case int:
		return util.ResultOk(fmt.Sprintf("'%d'", field.(int)))
	case string:
		return util.ResultOk("'" + field.(string) + "'")
	case bool:
		if field.(bool) {
			return util.ResultOk("'true'")
		}
		return util.ResultOk("'false'")
	case time.Time:
		return util.ResultOk("'" + field.(time.Time).Format(PostgresTimeFormat) + "'")
	}
	return util.ResultErr[string](fmt.Errorf("field: %+v has unsupported type", field))
}
