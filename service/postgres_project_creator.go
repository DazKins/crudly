package service

import (
	"crudly/model"
	"database/sql"
	"fmt"
)

type postgresProjectCreator struct {
	postgres *sql.DB
}

func NewPostgresProjectCreator(postgres *sql.DB) postgresProjectCreator {
	return postgresProjectCreator{
		postgres,
	}
}

func (p postgresProjectCreator) CreateProject(
	id model.ProjectId,
	authInfo model.ProjectAuthInfo,
) error {
	query := getPostgresProjectCreationQuery(id, authInfo.Salt, authInfo.SaltedHash)

	_, err := p.postgres.Query(query)

	if err != nil {
		return fmt.Errorf("error executing postgres query: %w", err)
	}

	return nil
}

func getPostgresProjectCreationQuery(id model.ProjectId, salt string, saltedHash string) string {
	return "INSERT INTO projects(id, salt, saltedhash) VALUES ('" + id.String() + "','" + salt + "','" + saltedHash + "')"
}
