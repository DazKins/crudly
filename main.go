package main

import (
	"crudly/app"
	"crudly/app/validation"
	"crudly/config"
	"crudly/http"
	"crudly/postgres"
	"crudly/redis"
	"crudly/service"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	config := config.InitialiseConfg()

	postgres, err := postgres.NewPostgres(config)

	if err != nil {
		fmt.Printf("error initialising postgres: %s", err.Error())
		return
	}

	redis := redis.NewRedis(config)

	postgresTableCreatorService := service.NewPostgresTableCreator(postgres)
	postgresTableGetterService := service.NewPostgresTableFetcher(postgres)
	postgresTableDeleterService := service.NewPostgresTableDeleter(postgres)
	postgresEntityFetcherService := service.NewPostgresEntityFetcher(postgres)
	postgresEntityCreatorService := service.NewPostgresEntityCreator(postgres)
	postgresEntityUpdaterService := service.NewPostgresEntityUpdater(postgres)
	postgresProjectCreatorService := service.NewPostgresProjectCreator(postgres)
	postgresProjectAuthInfoFetcherService := service.NewPostgresProjectAuthFetcher(postgres)
	postgresEntityDeleterService := service.NewPostgresEntityDeleter(postgres)

	redisRateLimitStoreService := service.NewRedisRateLimiterStore(redis)

	entityValidator := validation.NewEntityValidator()
	partialEntityValidator := validation.NewPartialEntityValidator()
	entityFilterValidator := validation.NewEntityFilterValidator()
	entityOrderValidator := validation.NewEntityOrderValidator()
	tableSchemaValidator := validation.NewTableSchemaValidator()

	projectManager := app.NewProjectManager(&postgresProjectCreatorService, &postgresProjectAuthInfoFetcherService)
	tableManager := app.NewTableManager(
		&postgresTableGetterService,
		&postgresTableCreatorService,
		&postgresTableDeleterService,
		&tableSchemaValidator,
	)
	entityManager := app.NewEntityManager(
		&postgresEntityFetcherService,
		&postgresEntityCreatorService,
		&postgresEntityUpdaterService,
		&postgresEntityDeleterService,
		&tableManager,
		&entityValidator,
		&partialEntityValidator,
		&entityFilterValidator,
		&entityOrderValidator,
	)
	rateLimitManager := app.NewRateLimitManager(&redisRateLimitStoreService)

	http.StartServer(
		config,
		&projectManager,
		&tableManager,
		&entityManager,
		&rateLimitManager,
	)
}
