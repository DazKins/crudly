package service

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"database/sql"
	"fmt"
)

type postgresRateLimitStore struct {
	postgres *sql.DB
}

func NewPostgresRateLimitStore(postgres *sql.DB) postgresRateLimitStore {
	return postgresRateLimitStore{
		postgres,
	}
}

func (p *postgresRateLimitStore) SetRateLimit(projectId model.ProjectId, rateLimit uint) error {
	query := getPostgresUpdateRateLimitQuery(projectId, rateLimit)

	_, err := p.postgres.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func (p *postgresRateLimitStore) GetRateLimit(projectId model.ProjectId) result.R[uint] {
	query := getPostgresRateLimitQuery(projectId)

	rows, err := p.postgres.Query(query)

	if err != nil {
		return result.Errf[uint]("error querying postgres: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return result.Err[uint](errs.RateLimitNotFoundError{})
	}

	rateLimit := uint(0)

	rows.Scan(&rateLimit)

	return result.Ok(rateLimit)
}

func getPostgresUpdateRateLimitQuery(projectId model.ProjectId, rateLimit uint) string {
	return fmt.Sprintf(
		`INSERT INTO rateLimit (projectId, rateLimit)
			VALUES ('%s', %d)
			ON CONFLICT (projectId)
			DO UPDATE SET rateLimit = EXCLUDED.rateLimit;`,
		projectId.String(),
		rateLimit,
	)
}

func getPostgresRateLimitQuery(projectId model.ProjectId) string {
	return fmt.Sprintf("SELECT rateLimit FROM rateLimit WHERE projectId = '%s'", projectId.String())
}
