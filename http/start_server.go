package http

import (
	"crudly/config"
	"crudly/http/handler"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type projectManager interface {
	CreateProject() result.Result[model.CreateProjectResponse]
	GetProjectAuthInfo(id model.ProjectId) result.Result[model.ProjectAuthInfo]
}

type tableManager interface {
	CreateTable(projectId model.ProjectId, name model.TableName, schema model.TableSchema) error
	GetTableSchema(projectId model.ProjectId, name model.TableName) result.Result[model.TableSchema]
}

type entityManager interface {
	GetEntity(projectId model.ProjectId, tableName model.TableName, id model.EntityId) result.Result[model.Entity]
	GetEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		entityFilter model.EntityFilter,
		paginationParams model.PaginationParams,
	) result.Result[model.Entities]
	CreateEntityWithId(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		entity model.Entity,
	) error
	CreateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		entity model.Entity,
	) error
	UpdateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		partialEntity model.PartialEntity,
	) error
	DeleteEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
	) error
}

func createHandler(
	config config.Config,
	projectManager projectManager,
	tableManager tableManager,
	entityManager entityManager,
) http.Handler {
	adminProjectHandler := handler.NewAdminProjectHandler(config, projectManager)
	adminTableHandler := handler.NewAdminTableHandler(
		tableManager,
		tableManager,
	)
	entityHandler := handler.NewEntityHandler(
		entityManager,
		entityManager,
		entityManager,
		entityManager,
	)

	tableNameMiddleware := middleware.NewTableName()
	projectIdMiddleware := middleware.NewProjectId()
	projectAuthMiddleware := middleware.NewProjectAuth(projectManager)
	loggerMiddleware := middleware.NewLogger(os.Stdout)

	router := mux.NewRouter()
	router.Use(loggerMiddleware)

	projectRouter := router.PathPrefix("/projects").Subrouter()

	projectRouter.HandleFunc(
		"",
		adminProjectHandler.PostProject,
	).Methods("POST")

	tableRouter := router.PathPrefix("/tables").Subrouter()
	tableRouter.Use(projectIdMiddleware)
	tableRouter.Use(projectAuthMiddleware)
	tableRouter.Use(tableNameMiddleware)

	tableRouter.HandleFunc(
		"/{tableName}",
		adminTableHandler.PutTable,
	).Methods("PUT")

	tableRouter.HandleFunc(
		"/{tableName}",
		adminTableHandler.GetTable,
	).Methods("GET")

	entityRouter := tableRouter.PathPrefix("/{tableName}/entities").Subrouter()

	entityRouter.HandleFunc(
		"/{id}",
		entityHandler.GetEntity,
	).Methods("GET")

	entityRouter.HandleFunc(
		"/{id}",
		entityHandler.PatchEntity,
	).Methods("PATCH")

	entityRouter.HandleFunc(
		"/{id}",
		entityHandler.DeleteEntity,
	).Methods("DELETE")

	entityRouter.HandleFunc(
		"",
		entityHandler.GetEntities,
	).Methods("GET")

	entityRouter.HandleFunc(
		"/{id}",
		entityHandler.PutEntity,
	).Methods("PUT")

	entityRouter.HandleFunc(
		"",
		entityHandler.PostEntity,
	).Methods("POST")

	return router
}

func StartServer(
	config config.Config,
	projectManager projectManager,
	tableManager tableManager,
	entityManager entityManager,
) {
	handler := createHandler(
		config,
		projectManager,
		tableManager,
		entityManager,
	)

	fmt.Printf("Starting server on port %d...\n", config.Port)

	http.ListenAndServe(
		fmt.Sprintf(":%d", config.Port),
		handler,
	)
}
