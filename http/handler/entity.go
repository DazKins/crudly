package handler

import (
	"crudly/http/dto"
	"crudly/model"
	"crudly/util"
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
	) util.Result[model.Entity]

	GetEntities(
		projectId model.ProjectId,
		tableName model.TableName,
		paginationParams model.PaginationParams,
	) util.Result[model.Entities]
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

type entityHandler struct {
	entityGetter  entityGetter
	entityCreator entityCreator
}

func NewEntityHandler(entityGetter entityGetter, entityCreator entityCreator) entityHandler {
	return entityHandler{
		entityGetter,
		entityCreator,
	}
}

func (e entityHandler) GetEntity(w http.ResponseWriter, r *http.Request) {
	projectId := util.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		w.WriteHeader(400)
		return
	}

	entityIdDto := dto.EntityIdDto(vars["id"])

	entityIdResult := entityIdDto.ToModel()

	if entityIdResult.IsErr() {
		w.WriteHeader(400)
		return
	}

	entityResult := e.entityGetter.GetEntity(
		projectId,
		tableNameResult.Unwrap(),
		entityIdResult.Unwrap(),
	)

	if entityResult.IsErr() {
		w.WriteHeader(500)
		return
	}

	entityDto := dto.GetEntityDto(entityResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entityDto)

	w.Write(resBodyBytes)
}

func (e entityHandler) GetEntities(w http.ResponseWriter, r *http.Request) {
	projectId := util.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		w.WriteHeader(400)
		return
	}

	paginationParams := model.PaginationParams{
		Limit:  model.DefaultPaginationLimit,
		Offset: model.DefaultPaginationOffset,
	}

	paginationLimitPathParam := dto.PaginationLimitPathParam(r.URL.Query().Get("limit"))

	if paginationLimitPathParam != "" {
		limitResult := paginationLimitPathParam.ToModel()

		if limitResult.IsErr() {
			w.WriteHeader(400)
			return
		}

		paginationParams.Limit = limitResult.Unwrap()
	}

	paginationOffsetPathParam := dto.PaginationOffsetPathParam(r.URL.Query().Get("offset"))

	if paginationOffsetPathParam != "" {
		offsetResult := paginationOffsetPathParam.ToModel()

		if offsetResult.IsErr() {
			w.WriteHeader(400)
			return
		}

		paginationParams.Offset = offsetResult.Unwrap()
	}

	entitiesResult := e.entityGetter.GetEntities(
		projectId,
		tableNameResult.Unwrap(),
		paginationParams,
	)

	if entitiesResult.IsErr() {
		w.WriteHeader(500)
		return
	}

	entitiesDto := dto.GetEntitiesDto(entitiesResult.Unwrap())

	resBodyBytes, _ := json.Marshal(entitiesDto)

	w.Write(resBodyBytes)
}

func (e entityHandler) PutEntity(w http.ResponseWriter, r *http.Request) {
	projectId := util.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		w.WriteHeader(400)
		return
	}

	entityIdDto := dto.EntityIdDto(vars["id"])

	entityIdResult := entityIdDto.ToModel()

	if entityIdResult.IsErr() {
		w.WriteHeader(400)
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
		w.WriteHeader(400)
		return
	}

	err = e.entityCreator.CreateEntityWithId(
		projectId,
		tableNameResult.Unwrap(),
		entityIdResult.Unwrap(),
		entityResult.Unwrap(),
	)

	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func (e entityHandler) PostEntity(w http.ResponseWriter, r *http.Request) {
	projectId := util.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		w.WriteHeader(400)
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
		w.WriteHeader(400)
		return
	}

	err = e.entityCreator.CreateEntity(
		projectId,
		tableNameResult.Unwrap(),
		entityResult.Unwrap(),
	)

	if err != nil {
		w.WriteHeader(500)
		return
	}
}
