package handler

import (
	"crudly/ctx"
	"crudly/errs"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util/result"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type entityGetter interface {
	GetEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
	) result.R[model.Entity]

	GetEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		entityFilter model.EntityFilter,
		entityOrder model.EntityOrder,
		paginationParams model.PaginationParams,
	) result.R[model.Entities]
}

type entityCreator interface {
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
}

type entityUpdater interface {
	UpdateEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
		partialEntity model.PartialEntity,
	) result.R[model.Entity]
}

type entityDeleter interface {
	DeleteEntity(
		projectId model.ProjectId,
		tableName model.TableName,
		id model.EntityId,
	) error
}

type entityHandler struct {
	entityGetter  entityGetter
	entityCreator entityCreator
	entityUpdater entityUpdater
	entityDeleter entityDeleter
}

func NewEntityHandler(
	entityGetter entityGetter,
	entityCreator entityCreator,
	entityUpdater entityUpdater,
	entityDeleter entityDeleter,
) entityHandler {
	return entityHandler{
		entityGetter,
		entityCreator,
		entityUpdater,
		entityDeleter,
	}
}

func (e *entityHandler) GetEntity(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	vars := mux.Vars(r)

	entityIdDto := dto.EntityIdDto(vars["id"])

	entityIdResult := entityIdDto.ToModel()

	if entityIdResult.IsErr() {
		middleware.AttachError(w, entityIdResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity id"))
		return
	}

	entityResult := e.entityGetter.GetEntity(
		projectId,
		tableName,
		entityIdResult.Unwrap(),
	)

	if entityResult.IsErr() {
		err := entityResult.UnwrapErr()

		middleware.AttachError(w, err)

		if _, ok := err.(errs.EntityNotFoundError); ok {
			w.WriteHeader(404)
			w.Write([]byte("entity not found"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting entity"))
		return
	}

	entityDto := dto.GetEntityDto(entityResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entityDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (e *entityHandler) GetEntities(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	paginationParams := model.PaginationParams{
		Limit:  model.DefaultPaginationLimit,
		Offset: model.DefaultPaginationOffset,
	}

	paginationLimitPathParam := dto.PaginationLimitPathParam(r.URL.Query().Get("limit"))

	if paginationLimitPathParam != "" {
		limitResult := paginationLimitPathParam.ToModel()

		if limitResult.IsErr() {
			middleware.AttachError(w, limitResult.UnwrapErr())
			w.WriteHeader(400)
			w.Write([]byte("invalid limit query param"))
			return
		}

		paginationParams.Limit = limitResult.Unwrap()
	}

	paginationOffsetPathParam := dto.PaginationOffsetPathParam(r.URL.Query().Get("offset"))

	if paginationOffsetPathParam != "" {
		offsetResult := paginationOffsetPathParam.ToModel()

		if offsetResult.IsErr() {
			middleware.AttachError(w, offsetResult.UnwrapErr())
			w.WriteHeader(400)
			w.Write([]byte("invalid offset query param"))
			return
		}

		paginationParams.Offset = offsetResult.Unwrap()
	}

	entityFilterResult := dto.GetEntityFilterFromQuery(r.URL.Query())

	if entityFilterResult.IsErr() {
		middleware.AttachError(w, entityFilterResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte(entityFilterResult.UnwrapErr().Error()))
		return
	}

	entityOrderResult := dto.GetEntityOrderFromQuery(r.URL.Query())

	if entityOrderResult.IsErr() {
		err := entityOrderResult.UnwrapErr()

		middleware.AttachError(w, err)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	entitiesResult := e.entityGetter.GetEntities(
		projectId,
		tableName,
		entityFilterResult.Unwrap(),
		entityOrderResult.Unwrap(),
		paginationParams,
	)

	if entitiesResult.IsErr() {
		err := entitiesResult.UnwrapErr()

		middleware.AttachError(w, err)

		if invalidEntityFilterError, ok := err.(errs.InvalidEntityFilterError); ok {
			w.WriteHeader(400)
			w.Write([]byte(invalidEntityFilterError.Error()))
			return
		}

		if invalidEntityOrderError, ok := err.(errs.InvalidEntityOrderError); ok {
			w.WriteHeader(400)
			w.Write([]byte(invalidEntityOrderError.Error()))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting entities"))
		return
	}

	entitiesDto := dto.GetEntitiesDto(entitiesResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entitiesDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (e *entityHandler) PutEntity(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	vars := mux.Vars(r)

	entityIdDto := dto.EntityIdDto(vars["id"])

	entityIdResult := entityIdDto.ToModel()

	if entityIdResult.IsErr() {
		middleware.AttachError(w, entityIdResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity id"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var entityDto dto.EntityDto
	json.Unmarshal(bodyBytes, &entityDto)

	entityResult := entityDto.ToModel()

	if entityResult.IsErr() {
		middleware.AttachError(w, entityResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity"))
		return
	}

	err = e.entityCreator.CreateEntityWithId(
		projectId,
		tableName,
		entityIdResult.Unwrap(),
		entityResult.Unwrap(),
	)

	if err != nil {
		middleware.AttachError(w, err)

		if err, ok := err.(errs.InvalidEntityError); ok {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		if _, ok := err.(errs.EntityAlreadyExistsError); ok {
			w.WriteHeader(409)
			w.Write([]byte("entity already exists"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating entity"))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(entityIdResult.Unwrap().String()))
}

func (e *entityHandler) PostEntity(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var entityDto dto.EntityDto
	json.Unmarshal(bodyBytes, &entityDto)

	entityResult := entityDto.ToModel()

	if entityResult.IsErr() {
		middleware.AttachError(w, entityResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity"))
		return
	}

	entityIdResult := e.entityCreator.CreateEntity(
		projectId,
		tableName,
		entityResult.Unwrap(),
	)

	if entityIdResult.IsErr() {
		err := entityIdResult.UnwrapErr()

		middleware.AttachError(w, err)

		if err, ok := err.(errs.InvalidEntityError); ok {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating entity"))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(entityIdResult.Unwrap().String()))
}

func (e *entityHandler) PostEntityBatch(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var entitiesDto dto.EntitiesDto
	json.Unmarshal(bodyBytes, &entitiesDto)

	entitiesResult := entitiesDto.ToModel()

	if entitiesResult.IsErr() {
		middleware.AttachError(w, entitiesResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entities array"))
		return
	}

	err = e.entityCreator.CreateEntities(
		projectId,
		tableName,
		entitiesResult.Unwrap(),
	)

	if err != nil {
		middleware.AttachError(w, err)

		if err, ok := err.(errs.InvalidEntityError); ok {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating entity"))
		return
	}

	w.WriteHeader(201)
}

func (e *entityHandler) PatchEntity(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	vars := mux.Vars(r)

	entityIdDto := dto.EntityIdDto(vars["id"])

	entityIdResult := entityIdDto.ToModel()

	if entityIdResult.IsErr() {
		middleware.AttachError(w, entityIdResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity id"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var partialEntityDto dto.PartialEntityDto
	json.Unmarshal(bodyBytes, &partialEntityDto)

	partialEntityResult := partialEntityDto.ToModel()

	if partialEntityResult.IsErr() {
		middleware.AttachError(w, partialEntityResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity"))
		return
	}

	entityResult := e.entityUpdater.UpdateEntity(
		projectId,
		tableName,
		entityIdResult.Unwrap(),
		partialEntityResult.Unwrap(),
	)

	if entityResult.IsErr() {
		err := entityResult.UnwrapErr()

		middleware.AttachError(w, err)

		if _, ok := err.(errs.InvalidPartialEntityError); ok {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error updating entity"))
		return
	}

	entityDto := dto.GetEntityDto(entityResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entityDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (e *entityHandler) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	vars := mux.Vars(r)

	entityIdDto := dto.EntityIdDto(vars["id"])

	entityIdResult := entityIdDto.ToModel()

	if entityIdResult.IsErr() {
		middleware.AttachError(w, entityIdResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid entity id"))
		return
	}

	err := e.entityDeleter.DeleteEntity(
		projectId,
		tableName,
		entityIdResult.Unwrap(),
	)

	if err != nil {
		middleware.AttachError(w, err)

		if _, ok := err.(errs.EntityNotFoundError); ok {
			w.WriteHeader(404)
			w.Write([]byte("entity not found"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error deleting entity"))
		return
	}
}
