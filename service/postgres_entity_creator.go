package service

import (
	"context"
	"crudly/errs"
	"crudly/model"
	"crudly/util"
	"database/sql"
	"fmt"

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

func (p *postgresEntityCreator) CreateEntity(
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

	rows, err := p.postgres.Query(query)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errs.EntityAlreadyExistsError{}
			}
		}

		return fmt.Errorf("error querying postgres: %w", err)
	}

	defer rows.Close()

	return nil
}

func (p *postgresEntityCreator) CreateEntities(
	projectId model.ProjectId,
	tableName model.TableName,
	ids []model.EntityId,
	entities model.Entities,
) error {
	tx, err := p.postgres.BeginTx(context.Background(), nil)

	if err != nil {
		return fmt.Errorf("error opening postgres transaction: %w", err)
	}
	defer tx.Rollback()

	for index, entity := range entities {
		query := getPostgresCreateEntityQuery(
			projectId,
			tableName,
			ids[index],
			entity,
		)

		_, err := tx.Query(query)

		if err != nil {
			return fmt.Errorf("error querying postgres: %w", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error commiting postgres transaction: %w", err)
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
		query = query + "\"" + k.String() + "\","
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
