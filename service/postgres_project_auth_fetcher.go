package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util"
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

func (p postgresProjectAuthFetcher) FetchProjectAuthInfo(id model.ProjectId) util.Result[model.ProjectAuthInfo] {
	query := getPostgresFetchProjectAuthQuery(id)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return util.ResultErr[model.ProjectAuthInfo](fmt.Errorf("error querying postgres: %w", err))
	}

	if !rows.Next() {
		return util.ResultErr[model.ProjectAuthInfo](errs.ProjectNotFoundError{})
	}

	salt := ""
	saltedHash := ""

	rows.Scan(&salt, &saltedHash)

	return util.ResultOk(model.ProjectAuthInfo{
		Salt:       salt,
		SaltedHash: saltedHash,
	})
}

func getPostgresFetchProjectAuthQuery(id model.ProjectId) string {
	return "SELECT salt, saltedhash FROM projects WHERE id = '" + id.String() + "'"
}
