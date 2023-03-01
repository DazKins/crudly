package app

import (
	"crudly/model"
	"crudly/util"
	"fmt"

	"github.com/google/uuid"
)

type entityFetcher interface {
	FetchEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		tableSchema model.TableSchema,
		id model.EntityId,
	) util.Result[model.Entity]

	FetchEntities(
		projectId model.ProjectId,
		table model.TableName,
		tableSchema model.TableSchema,
		paginationParams model.PaginationParams,
	) util.Result[model.Entities]
}

type entityCreator interface {
	CreateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		tableSchema model.TableSchema,
		id model.EntityId,
		entity model.Entity,
	) error
}

type tableSchemaGetter interface {
	GetTableSchema(projectId model.ProjectId, name model.TableName) util.Result[model.TableSchema]
}

type entityManager struct {
	entityFetcher     entityFetcher
	entityCreator     entityCreator
	tableSchemaGetter tableSchemaGetter
}

func NewEntityManager(
	entityFetcher entityFetcher,
	entityCreator entityCreator,
	tableSchemaGetter tableSchemaGetter,
) entityManager {
	return entityManager{
		entityFetcher,
		entityCreator,
		tableSchemaGetter,
	}
}

func (e entityManager) GetEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) util.Result[model.Entity] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return util.ResultErr[model.Entity](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
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
) util.Result[model.Entities] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return util.ResultErr[model.Entities](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
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

	return e.entityCreator.CreateEntity(
		projectId,
		tableName,
		tableSchemaResult.Unwrap(),
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
