package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/optional"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type postgresUserCreator struct {
	postgres *sql.DB
}

func NewPostgresUserCreator(postgres *sql.DB) postgresUserCreator {
	return postgresUserCreator{
		postgres,
	}
}

func (p postgresUserCreator) CreateUser(id model.UserId, user model.User) error {
	query := fmt.Sprintf(
		"INSERT INTO users(id, twitterid, googleid) VALUES ('%s', %s, %s)",
		id.String(),
		getPostgresOptionalStringValue(user.TwitterId),
		getPostgresOptionalStringValue(user.GoogleId),
	)

	_, err := p.postgres.Query(query)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errs.UserAlreadyExistsError{}
			}
		}

		return fmt.Errorf("error querying postgres: %w", err)
	}

	return nil
}

func getPostgresOptionalStringValue(str optional.O[string]) string {
	if str.IsNone() {
		return "NULL"
	} else {
		return fmt.Sprintf("'%s'", str.Unwrap())
	}
}
