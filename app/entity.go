package app

import (
	"crudly/errs"
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
	) result.R[model.Entity]

	FetchEntities(
		projectId model.ProjectId,
		table model.TableName,
		tableSchema model.TableSchema,
		entityFilter model.EntityFilter,
		entityOrder model.EntityOrder,
		paginationParams model.PaginationParams,
	) result.R[model.Entities]
}

type entityCreator interface {
	CreateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		entity model.Entity,
	) error

	CreateEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		ids []model.EntityId,
		entities model.Entities,
	) error
}

type entityUpdater interface {
	UpdateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		partialEntity model.PartialEntity,
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
	GetTableSchema(projectId model.ProjectId, name model.TableName) result.R[model.TableSchema]
}

type entityValidator interface {
	ValidateEntity(entity model.Entity, tableSchema model.TableSchema) error
}

type partialEntityValidator interface {
	ValidatePartialEntity(partialEntity model.PartialEntity, tableSchema model.TableSchema) error
}

type entityFilterValidator interface {
	ValidateEntityFilter(
		entityFilter model.EntityFilter,
		tableSchema model.TableSchema,
	) error
}

type entityOrderValidator interface {
	ValidateEntityOrder(
		entityOrder model.EntityOrder,
		tableSchema model.TableSchema,
	) error
}

type entityManager struct {
	entityFetcher          entityFetcher
	entityCreator          entityCreator
	entityUpdater          entityUpdater
	entityDeleter          entityDeleter
	tableSchemaGetter      tableSchemaGetter
	entityValidator        entityValidator
	partialEntityValidator partialEntityValidator
	entityFilterValidator  entityFilterValidator
	entityOrderValidator   entityOrderValidator
}

func NewEntityManager(
	entityFetcher entityFetcher,
	entityCreator entityCreator,
	entityUpdater entityUpdater,
	entityDeleter entityDeleter,
	tableSchemaGetter tableSchemaGetter,
	entityValidator entityValidator,
	partialEntityValidator partialEntityValidator,
	entityFilterValidator entityFilterValidator,
	entityOrderValidator entityOrderValidator,
) entityManager {
	return entityManager{
		entityFetcher,
		entityCreator,
		entityUpdater,
		entityDeleter,
		tableSchemaGetter,
		entityValidator,
		partialEntityValidator,
		entityFilterValidator,
		entityOrderValidator,
	}
}

func (e entityManager) GetEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) result.R[model.Entity] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return result.Err[model.Entity](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	entityResult := e.entityFetcher.FetchEntity(
		projectId,
		tableName,
		tableSchemaResult.Unwrap(),
		id,
	)

	if entityResult.IsErr() {
		err := entityResult.UnwrapErr()

		if _, ok := err.(errs.EntityNotFoundError); ok {
			return entityResult
		}
	}

	return entityResult
}

func (e entityManager) GetEntities(
	projectId model.ProjectId,
	tableName model.TableName,
	entityFilter model.EntityFilter,
	entityOrder model.EntityOrder,
	paginationParams model.PaginationParams,
) result.R[model.Entities] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return result.Err[model.Entities](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	tableSchema := tableSchemaResult.Unwrap()

	err := e.entityFilterValidator.ValidateEntityFilter(entityFilter, tableSchema)

	if err != nil {
		return result.Err[model.Entities](errs.NewInvalidEntityFilterError(err))
	}

	err = e.entityOrderValidator.ValidateEntityOrder(entityOrder, tableSchema)

	if err != nil {
		return result.Err[model.Entities](errs.NewInvalidEntityOrderError(err))
	}

	return e.entityFetcher.FetchEntities(
		projectId,
		tableName,
		tableSchema,
		entityFilter,
		entityOrder,
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
		return errs.NewInvalidEntityError(err)
	}

	err = e.entityCreator.CreateEntity(
		projectId,
		tableName,
		id,
		entity,
	)

	if err != nil {
		if _, ok := err.(errs.EntityAlreadyExistsError); ok {
			return err
		}

		return fmt.Errorf("error creating entity: %w", err)
	}

	return nil
}

func (e entityManager) CreateEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	entity model.Entity,
) result.R[model.EntityId] {
	id := model.EntityId(uuid.New())

	err := e.CreateEntityWithId(
		projectId,
		tableName,
		id,
		entity,
	)

	if err != nil {
		return result.Err[model.EntityId](err)
	}

	return result.Ok(id)
}

func (e entityManager) CreateEntities(
	projectId model.ProjectId,
	tableName model.TableName,
	entities model.Entities,
) error {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr())
	}

	tableSchema := tableSchemaResult.Unwrap()

	for index, entity := range entities {
		err := e.entityValidator.ValidateEntity(entity, tableSchema)

		if err != nil {
			return errs.NewInvalidEntityError(
				fmt.Errorf("error with entity at index %d: %w", index, err),
			)
		}
	}

	entityIds := make([]model.EntityId, len(entities))

	for index := range entities {
		entityIds[index] = model.EntityId(uuid.New())
	}

	err := e.entityCreator.CreateEntities(projectId, tableName, entityIds, entities)

	if err != nil {
		return fmt.Errorf("error creating entities: %w", err)
	}

	return nil
}

func (e entityManager) UpdateEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
	partialEntity model.PartialEntity,
) error {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr())
	}

	tableSchema := tableSchemaResult.Unwrap()

	err := e.partialEntityValidator.ValidatePartialEntity(partialEntity, tableSchema)

	if err != nil {
		return errs.NewInvalidPartialEntityError(err)
	}

	return e.entityUpdater.UpdateEntity(
		projectId,
		tableName,
		id,
		partialEntity,
	)
}

func (e entityManager) DeleteEntity(
	projectId model.ProjectId,
	tableName model.TableName,
	id model.EntityId,
) error {
	err := e.entityDeleter.DeleteEntity(
		projectId,
		tableName,
		id,
	)

	if err != nil {
		if _, ok := err.(errs.EntityNotFoundError); ok {
			return err
		}

		return fmt.Errorf("error deleting entity: %w", err)
	}

	return nil
}
