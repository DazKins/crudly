package app

import (
	"crudly/model"
	"crudly/util/result"
	"fmt"

	"github.com/google/uuid"
)

type entityFetcher interface {
	FetchEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		tableSchema model.TableSchema,
		id model.EntityId,
	) result.Result[model.Entity]

	FetchEntities(
		projectId model.ProjectId,
		table model.TableName,
		tableSchema model.TableSchema,
		paginationParams model.PaginationParams,
	) result.Result[model.Entities]
}

type entityCreator interface {
	CreateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		entity model.Entity,
	) error
}

type entityDeleter interface {
	DeleteEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
	) error
}

type tableSchemaGetter interface {
	GetTableSchema(projectId model.ProjectId, name model.TableName) result.Result[model.TableSchema]
}

type entityValidator interface {
	ValidateEntity(entity model.Entity, tableSchema model.TableSchema) error
}

type entityManager struct {
	entityFetcher     entityFetcher
	entityCreator     entityCreator
	entityDeleter     entityDeleter
	tableSchemaGetter tableSchemaGetter
	entityValidator   entityValidator
}

func NewEntityManager(
	entityFetcher entityFetcher,
	entityCreator entityCreator,
	entityDeleter entityDeleter,
	tableSchemaGetter tableSchemaGetter,
	entityValidator entityValidator,
) entityManager {
	return entityManager{
		entityFetcher,
		entityCreator,
		entityDeleter,
		tableSchemaGetter,
		entityValidator,
	}
}

func (e entityManager) GetEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) result.Result[model.Entity] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return result.Err[model.Entity](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	return e.entityFetcher.FetchEntity(
		projectId,
		tableName,
		tableSchemaResult.Unwrap(),
		id,
	)
}

func (e entityManager) GetEntities(
	projectId model.ProjectId,
	tableName model.TableName,
	paginationParams model.PaginationParams,
) result.Result[model.Entities] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return result.Err[model.Entities](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	return e.entityFetcher.FetchEntities(
		projectId,
		tableName,
		tableSchemaResult.Unwrap(),
		paginationParams,
	)
}

func (e entityManager) CreateEntityWithId(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
	entity model.Entity,
) error {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr())
	}

	tableSchema := tableSchemaResult.Unwrap()

	err := e.entityValidator.ValidateEntity(entity, tableSchema)

	if err != nil {
		return fmt.Errorf("error validating entity: %w", err)
	}

	return e.entityCreator.CreateEntity(
		projectId,
		tableName,
		id,
		entity,
	)
}

func (e entityManager) CreateEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	entity model.Entity,
) error {
	id := model.EntityId(uuid.New())

	return e.CreateEntityWithId(
		projectId,
		tableName,
		id,
		entity,
	)
}

func (e entityManager) DeleteEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) error {
	return e.entityDeleter.DeleteEntity(
		projectId,
		tableName,
		id,
	)
}
