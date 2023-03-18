package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util"
	"crudly/util/result"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errs.EntityAlreadyExistsError{}
			}
		}

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
		query = query + k.String() + ","
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

func getPostgresFieldValue(field any) result.Result[string] {
	switch v := field.(type) {
	case uuid.UUID:
		return result.Ok("'" + v.String() + "'")
	case int:
		return result.Ok(fmt.Sprintf("'%d'", v))
	case string:
		return result.Ok("'" + v + "'")
	case bool:
		if v {
			return result.Ok("'true'")
		}
		return result.Ok("'false'")
	case time.Time:
		return result.Ok("'" + v.Format(PostgresTimeFormat) + "'")
	}
	return result.Err[string](fmt.Errorf("field: %+v has unsupported type", field))
}
