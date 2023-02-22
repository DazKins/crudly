package main

import (
	"crudly/app"
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
	postgresEntityFetcherService := service.NewPostgresEntityFetcher(postgres)
	postgresEntityCreatorService := service.NewPostgresEntityCreator(postgres)
	postgresProjectCreatorService := service.NewPostgresProjectCreator(postgres)
	postgresProjectAuthInfoFetcherService := service.NewPostgresProjectAuthFetcher(postgres)

	projectManager := app.NewProjectManager(postgresProjectCreatorService, postgresProjectAuthInfoFetcherService)
	entityManager := app.NewEntityManager(postgresEntityFetcherService, postgresEntityCreatorService)
	tableManager := app.NewTableManager(postgresTableGetterService, postgresTableCreatorService)

	http.StartServer(
		config,
		projectManager,
		tableManager,
		entityManager,
	)
}
