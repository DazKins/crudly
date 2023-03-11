package app

import (
	"crudly/errs"
	"crudly/model"
	"crudly/util/result"
	"errors"
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
		entityFilter model.EntityFilter,
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

type entityFilterValidator interface {
	ValidateEntityFilter(
		entityFilter model.EntityFilter,
		tableSchema model.TableSchema,
	) error
}

type entityManager struct {
	entityFetcher         entityFetcher
	entityCreator         entityCreator
	entityDeleter         entityDeleter
	tableSchemaGetter     tableSchemaGetter
	entityValidator       entityValidator
	entityFilterValidator entityFilterValidator
}

func NewEntityManager(
	entityFetcher entityFetcher,
	entityCreator entityCreator,
	entityDeleter entityDeleter,
	tableSchemaGetter tableSchemaGetter,
	entityValidator entityValidator,
	entityFilterValidator entityFilterValidator,
) entityManager {
	return entityManager{
		entityFetcher,
		entityCreator,
		entityDeleter,
		tableSchemaGetter,
		entityValidator,
		entityFilterValidator,
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
	paginationParams model.PaginationParams,
) result.Result[model.Entities] {
	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableName)

	if tableSchemaResult.IsErr() {
		return result.Err[model.Entities](fmt.Errorf("error getting table schema: %w", tableSchemaResult.UnwrapErr()))
	}

	tableSchema := tableSchemaResult.Unwrap()

	err := e.entityFilterValidator.ValidateEntityFilter(entityFilter, tableSchema)

	if err != nil {
		return result.Err[model.Entities](errs.NewInvalidEntityFilterError(err))
	}

	return e.entityFetcher.FetchEntities(
		projectId,
		tableName,
		tableSchema,
		entityFilter,
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
	err := e.entityDeleter.DeleteEntity(
		projectId,
		tableName,
		id,
	)

	if err != nil {
		if errors.As(err, new(errs.EntityNotFoundError)) {
			return err
		}

		return fmt.Errorf("error deleting entity: %w", err)
	}

	return nil
}
