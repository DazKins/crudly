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
	) result.Result[model.Entity]

	GetEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		paginationParams model.PaginationParams,
	) result.Result[model.Entities]
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
	) error
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
	entityDeleter entityDeleter
}

func NewEntityHandler(
	entityGetter entityGetter,
	entityCreator entityCreator,
	entityDeleter entityDeleter,
) entityHandler {
	return entityHandler{
		entityGetter,
		entityCreator,
		entityDeleter,
	}
}

func (e entityHandler) GetEntity(w http.ResponseWriter, r *http.Request) {
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
		middleware.AttachError(w, entityResult.UnwrapErr())
		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting entity"))
		return
	}

	entityDto := dto.GetEntityDto(entityResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entityDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (e entityHandler) GetEntities(w http.ResponseWriter, r *http.Request) {
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

	entitiesResult := e.entityGetter.GetEntities(
		projectId,
		tableName,
		paginationParams,
	)

	if entitiesResult.IsErr() {
		err := entitiesResult.UnwrapErr()

		middleware.AttachError(w, err)

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting entities"))
		return
	}

	entitiesDto := dto.GetEntitiesDto(entitiesResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entitiesDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (e entityHandler) PutEntity(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating entity"))
		return
	}
}

func (e entityHandler) PostEntity(w http.ResponseWriter, r *http.Request) {
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

	err = e.entityCreator.CreateEntity(
		projectId,
		tableName,
		entityResult.Unwrap(),
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
}

func (e entityHandler) DeleteEntity(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(500)
		w.Write([]byte("unexpected error deleting entity"))
		return
	}
}
