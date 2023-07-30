package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"database/sql"
	"fmt"
)

type postgresEntityCount struct {
	postgres *sql.DB
}

func NewPostgresEntityCount(postgres *sql.DB) postgresEntityCount {
	return postgresEntityCount{
		postgres,
	}
}

func (p *postgresEntityCount) FetchTotalEntityCount(
	projectId model.ProjectId,
	tableName model.TableName,
) result.R[uint] {
	query := getPostgresRowCountQuery(projectId, tableName)

	rows, err := p.postgres.Query(query)

	if err != nil {
		result.Errf[uint]("error querying postgres: %w", err)
	}

	if !rows.Next() {
		return result.Err[uint](errs.TableNotFoundError{})
	}

	totalCount := uint(0)

	rows.Scan(&totalCount)

	return result.Ok(totalCount)
}

func getPostgresRowCountQuery(
	projectId model.ProjectId,
	tableName model.TableName,
) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", getPostgresTableName(projectId, tableName))
}
