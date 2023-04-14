package main

import (
	"crudly/app"
	"crudly/app/validation"
	"crudly/config"
	"crudly/http"
	"crudly/postgres"
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

	postgresTableCreatorService := service.NewPostgresTableCreator(postgres)
	postgresTableGetterService := service.NewPostgresTableFetcher(postgres)
	postgresTableDeleterService := service.NewPostgresTableDeleter(postgres)
	postgresEntityFetcherService := service.NewPostgresEntityFetcher(postgres)
	postgresEntityCreatorService := service.NewPostgresEntityCreator(postgres)
	postgresEntityUpdaterService := service.NewPostgresEntityUpdater(postgres)
	postgresProjectCreatorService := service.NewPostgresProjectCreator(postgres)
	postgresProjectAuthInfoFetcherService := service.NewPostgresProjectAuthFetcher(postgres)
	postgresEntityDeleterService := service.NewPostgresEntityDeleter(postgres)
	postgresUserFetcherService := service.NewPostgresUserFetcher(postgres)
	postgresUserCreatorService := service.NewPostgresUserCreator(postgres)

	entityValidator := validation.NewEntityValidator()
	partialEntityValidator := validation.NewPartialEntityValidator()
	entityFilterValidator := validation.NewEntityFilterValidator()
	entityOrderValidator := validation.NewEntityOrderValidator()
	tableSchemaValidator := validation.NewTableSchemaValidator()

	userManager := app.NewUserManager(postgresUserCreatorService, postgresUserFetcherService)
	projectManager := app.NewProjectManager(postgresProjectCreatorService, postgresProjectAuthInfoFetcherService)
	tableManager := app.NewTableManager(
		postgresTableGetterService,
		postgresTableCreatorService,
		postgresTableDeleterService,
		tableSchemaValidator,
	)
	entityManager := app.NewEntityManager(
		postgresEntityFetcherService,
		postgresEntityCreatorService,
		postgresEntityUpdaterService,
		postgresEntityDeleterService,
		tableManager,
		entityValidator,
		partialEntityValidator,
		entityFilterValidator,
		entityOrderValidator,
	)

	http.StartServer(
		config,
		userManager,
		projectManager,
		tableManager,
		entityManager,
	)
}
