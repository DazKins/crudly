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
	GetEntities(projectId model.ProjectId, tableName model.TableName, paginationParams model.PaginationParams) result.Result[model.Entities]
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
	entityHandler := handler.NewEntityHandler(entityManager, entityManager)

	router := mux.NewRouter()

	adminRouter := router.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/project", adminProjectHandler.PostProject).Methods("POST")

	projectIdMiddleware := middleware.NewProjectId()
	projectAuthMiddleware := middleware.NewProjectAuth(projectManager)
	loggerMiddleware := middleware.NewLogger(os.Stdout)

	projectAuthMiddlewares := middleware.Middlewares{
		projectIdMiddleware,
		projectAuthMiddleware,
	}

	adminRouter.HandleFunc(
		"/table/{tableName}",
		middleware.AttachMultiple(
			adminTableHandler.PutTable,
			projectAuthMiddlewares,
		),
	).Methods("PUT")

	adminRouter.HandleFunc(
		"/table/{tableName}",
		middleware.AttachMultiple(
			adminTableHandler.GetTable,
			projectAuthMiddlewares,
		),
	).Methods("GET")

	router.HandleFunc(
		"/{tableName}/{id}",
		middleware.AttachMultiple(
			entityHandler.GetEntity,
			projectAuthMiddlewares,
		),
	).Methods("GET")

	router.HandleFunc(
		"/{tableName}",
		middleware.AttachMultiple(
			entityHandler.GetEntities,
			projectAuthMiddlewares,
		),
	).Methods("GET")

	router.HandleFunc(
		"/{tableName}/{id}",
		middleware.AttachMultiple(
			entityHandler.PutEntity,
			projectAuthMiddlewares,
		),
	).Methods("PUT")

	router.HandleFunc(
		"/{tableName}",
		middleware.AttachMultiple(
			entityHandler.PostEntity,
			projectAuthMiddlewares,
		),
	).Methods("POST")

	loggedHandler := middleware.AttachToHandler(
		router,
		loggerMiddleware,
	)

	return loggedHandler
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
