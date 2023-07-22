package postgres

import (
	"crudly/config"
	"database/sql"
	"fmt"
)

func NewPostgres(config config.Config) (*sql.DB, error) {
	return sql.Open(
		"postgres",
		fmt.Sprintf(
			"sslmode=%s dbname=%s host=%s port=%d user=%s password=%s",
			config.PostgresSslMode,
			config.PostgresDatabase,
			config.PostgresHost,
			config.PostgresPort,
			config.PostgresUsername,
			config.PostgresPassword,
		),
	)
}
