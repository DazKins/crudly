package app

import (
	"crudly/model"
	"crudly/util"

	"github.com/google/uuid"
)

type entityFetcher interface {
	FetchEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
	) util.Result[model.Entity]

	FetchEntities(
		projectId model.ProjectId,
		table model.TableName,
		paginationParams model.PaginationParams,
	) util.Result[model.Entities]
}

type entityCreator interface {
	CreateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		entity model.Entity,
	) error
}

type entityManager struct {
	entityFetcher entityFetcher
	entityCreator entityCreator
}

func NewEntityManager(entityFetcher entityFetcher, entityCreator entityCreator) entityManager {
	return entityManager{
		entityFetcher,
		entityCreator,
	}
}

func (e entityManager) GetEntity(projectId model.ProjectId, tableName model.TableName, id model.EntityId) util.Result[model.Entity] {
	return e.entityFetcher.FetchEntity(projectId, tableName, id)
}

func (e entityManager) GetEntities(projectId model.ProjectId, tableName model.TableName, paginationParams model.PaginationParams) util.Result[model.Entities] {
	return e.entityFetcher.FetchEntities(projectId, tableName, paginationParams)
}

func (e entityManager) CreateEntityWithId(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
	entity model.Entity,
) error {
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
