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
	projectCreationQuery := getPostgresProjectCreationQuery(id, authInfo.Salt, authInfo.SaltedHash)

	_, err := p.postgres.Exec(projectCreationQuery)

	if err != nil {
		return fmt.Errorf("error executing postgres query: %w", err)
	}

	schemaTableCreationQuery := getPostgresSchemaTableCreationQuery(id)

	_, err = p.postgres.Exec(schemaTableCreationQuery)

	if err != nil {
		return fmt.Errorf("error executing postgres query: %w", err)
	}

	return nil
}

func getPostgresSchemaTableCreationQuery(id model.ProjectId) string {
	return "CREATE TABLE \"" + getPostgresSchemaTableName(id) + "\"(name varchar, schema varchar)"
}

func getPostgresProjectCreationQuery(id model.ProjectId, salt string, saltedHash string) string {
	return "INSERT INTO projects(id, salt, saltedhash) VALUES ('" + id.String() + "','" + salt + "','" + saltedHash + "')"
}
