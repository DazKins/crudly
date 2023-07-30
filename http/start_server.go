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
	CreateProject() result.R[model.CreateProjectResponse]
	GetProjectAuthInfo(id model.ProjectId) result.R[model.ProjectAuthInfo]
}

type tableManager interface {
	CreateTable(projectId model.ProjectId, name model.TableName, schema model.TableSchema) error
	GetTableSchema(projectId model.ProjectId, name model.TableName) result.R[model.TableSchema]
	GetTableSchemas(projectId model.ProjectId) result.R[model.TableSchemas]
	DeleteTable(projectId model.ProjectId, name model.TableName) error
}

type entityManager interface {
	GetEntity(projectId model.ProjectId, tableName model.TableName, id model.EntityId) result.R[model.Entity]
	GetEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		entityFilter model.EntityFilter,
		entityOrders model.EntityOrders,
		paginationParams model.PaginationParams,
	) result.R[model.Entities]
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
	) result.R[model.EntityId]
	CreateEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		entities model.Entities,
	) error
	UpdateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		partialEntity model.PartialEntity,
	) result.R[model.Entity]
	DeleteEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
	) error
}

type rateLimitManager interface {
	GetDailyRateLimit(projectId model.ProjectId) result.R[uint]
	GetCurrentRateUsage(projectId model.ProjectId) result.R[uint]
	HandleUsage(projectId model.ProjectId) error
	ShouldBlockRequest(projectId model.ProjectId) bool
}

func createHandler(
	config config.Config,
	projectManager projectManager,
	tableManager tableManager,
	entityManager entityManager,
	rateLimitManager rateLimitManager,
) http.Handler {
	projectHandler := handler.NewProjectHandler(config, projectManager)
	tableHandler := handler.NewTableHandler(
		tableManager,
		tableManager,
		tableManager,
	)
	entityHandler := handler.NewEntityHandler(
		entityManager,
		entityManager,
		entityManager,
		entityManager,
	)
	rateLimitHandler := handler.NewRateLimitHandler(rateLimitManager)

	adminApiKeyMiddleware := middleware.NewAdminApiKey(config)
	adminAuthMiddleware := middleware.NewAdminAuth()
	tableNameMiddleware := middleware.NewTableName()
	projectIdMiddleware := middleware.NewProjectId()
	projectAuthMiddleware := middleware.NewProjectAuth(projectManager)
	loggerMiddleware := middleware.NewLogger(os.Stdout)
	rateLimitMiddleware := middleware.NewRateLimit(rateLimitManager)

	router := mux.NewRouter()
	router.Use(loggerMiddleware)
	router.Use(adminApiKeyMiddleware)

	projectRouter := router.PathPrefix("/projects").Subrouter()
	projectRouter.Use(adminAuthMiddleware)

	projectRouter.HandleFunc(
		"",
		projectHandler.PostProject,
	).Methods("POST")

	rateLimitRouter := router.PathPrefix("/rateLimit").Subrouter()
	rateLimitRouter.Use(projectIdMiddleware)
	rateLimitRouter.Use(projectAuthMiddleware)

	rateLimitRouter.HandleFunc(
		"",
		rateLimitHandler.GetRateLimit,
	).Methods("GET")

	tableRouter := router.PathPrefix("/tables").Subrouter()
	tableRouter.Use(projectIdMiddleware)
	tableRouter.Use(projectAuthMiddleware)
	tableRouter.Use(rateLimitMiddleware)
	tableRouter.Use(tableNameMiddleware)

	tableRouter.HandleFunc(
		"/{tableName}",
		tableHandler.PutTable,
	).Methods("PUT")

	tableRouter.HandleFunc(
		"",
		tableHandler.GetTables,
	).Methods("GET")

	tableRouter.HandleFunc(
		"/{tableName}",
		tableHandler.GetTable,
	).Methods("GET")

	tableRouter.HandleFunc(
		"/{tableName}",
		tableHandler.DeleteTable,
	).Methods("DELETE")

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

	entityRouter.HandleFunc(
		"/batch",
		entityHandler.PostEntityBatch,
	).Methods("POST")

	return router
}

func StartServer(
	config config.Config,
	projectManager projectManager,
	tableManager tableManager,
	entityManager entityManager,
	rateLimitManager rateLimitManager,
) {
	handler := createHandler(
		config,
		projectManager,
		tableManager,
		entityManager,
		rateLimitManager,
	)

	fmt.Printf("Starting server on port %d...\n", config.Port)

	http.ListenAndServe(
		fmt.Sprintf(":%d", config.Port),
		handler,
	)
}
