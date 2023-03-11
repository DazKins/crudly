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

	router.HandleFunc(
		"/project",
		adminProjectHandler.PostProject,
	).Methods("POST")

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(projectIdMiddleware)
	adminRouter.Use(projectAuthMiddleware)

	adminRouter.HandleFunc(
		"/table/{tableName}",
		adminTableHandler.PutTable,
	).Methods("PUT")

	adminRouter.HandleFunc(
		"/table/{tableName}",
		adminTableHandler.GetTable,
	).Methods("GET")

	userRouter := router.PathPrefix("/{tableName}").Subrouter()
	userRouter.Use(projectIdMiddleware)
	userRouter.Use(projectAuthMiddleware)
	userRouter.Use(tableNameMiddleware)

	userRouter.HandleFunc(
		"/{id}",
		entityHandler.GetEntity,
	).Methods("GET")

	userRouter.HandleFunc(
		"/{id}",
		entityHandler.PatchEntity,
	).Methods("PATCH")

	userRouter.HandleFunc(
		"/{id}",
		entityHandler.DeleteEntity,
	).Methods("DELETE")

	userRouter.HandleFunc(
		"",
		entityHandler.GetEntities,
	).Methods("GET")

	userRouter.HandleFunc(
		"/{id}",
		entityHandler.PutEntity,
	).Methods("PUT")

	userRouter.HandleFunc(
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
