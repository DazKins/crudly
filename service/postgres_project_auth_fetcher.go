package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"database/sql"
	"fmt"
)

type postgresProjectAuthFetcher struct {
	postgres *sql.DB
}

func NewPostgresProjectAuthFetcher(postgres *sql.DB) postgresProjectAuthFetcher {
	return postgresProjectAuthFetcher{
		postgres,
	}
}

func (p postgresProjectAuthFetcher) FetchProjectAuthInfo(id model.ProjectId) result.Result[model.ProjectAuthInfo] {
	query := getPostgresFetchProjectAuthQuery(id)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.ProjectAuthInfo](fmt.Errorf("error querying postgres: %w", err))
	}

	if !rows.Next() {
		return result.Err[model.ProjectAuthInfo](errs.ProjectNotFoundError{})
	}

	salt := ""
	saltedHash := ""

	rows.Scan(&salt, &saltedHash)

	return result.Ok(model.ProjectAuthInfo{
		Salt:       salt,
		SaltedHash: saltedHash,
	})
}

func getPostgresFetchProjectAuthQuery(id model.ProjectId) string {
	return "SELECT salt, saltedhash FROM projects WHERE id = '" + id.String() + "'"
}
