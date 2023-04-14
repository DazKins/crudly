package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/optional"
	"crudly/util/result"
	"database/sql"
	"fmt"
)

type postgresUserFetcher struct {
	postgres *sql.DB
}

func NewPostgresUserFetcher(postgres *sql.DB) postgresUserFetcher {
	return postgresUserFetcher{
		postgres,
	}
}

func (p postgresUserFetcher) FetchUser(id model.UserId) result.R[model.User] {
	query := fmt.Sprintf(
		"SELECT %s FROM users WHERE id = '%s'",
		"twitterid, googleid",
		id.String(),
	)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Err[model.User](fmt.Errorf("error querying postgres: %w", err))
	}

	defer rows.Close()

	if !rows.Next() {
		return result.Err[model.User](errs.UserNotFoundError{})
	}

	twitterid := sql.NullString{}
	googleid := sql.NullString{}

	rows.Scan(&twitterid, &googleid)

	user := model.User{}

	if twitterid.Valid {
		user.TwitterId = optional.Some(twitterid.String)
	}

	if googleid.Valid {
		user.GoogleId = optional.Some(googleid.String)
	}

	return result.Ok(user)
}
