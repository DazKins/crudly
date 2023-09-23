package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"crudly/ctx"
	"crudly/errs"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util/optional"
	"crudly/util/result"

	"github.com/gorilla/mux"
)

type tableCreator interface {
	CreateTable(
		projectId model.ProjectId,
		tableName model.TableName,
		schema model.TableSchema,
	) error
}

type tableSchemaGetter interface {
	GetTableSchema(projectId model.ProjectId, name model.TableName) result.R[model.TableSchema]
	GetTableSchemas(projectId model.ProjectId) result.R[model.TableSchemas]
}

type tableDeleter interface {
	DeleteTable(projectId model.ProjectId, name model.TableName) error
}

type tableFieldAdder interface {
	AddField(
		projectId model.ProjectId,
		tableName model.TableName,
		name model.FieldName,
		definition model.FieldDefinition,
		defaultValue optional.O[any],
	) error
}

type tableFieldDeleter interface {
	DeleteField(
		projectId model.ProjectId,
		tableName model.TableName,
		name model.FieldName,
	) error
}

type tableHandler struct {
	tableCreator      tableCreator
	tableSchemaGetter tableSchemaGetter
	tableDeleter      tableDeleter
	tableFieldAdder   tableFieldAdder
	tableFieldDeleter tableFieldDeleter
}

func NewTableHandler(
	tableCreator tableCreator,
	tableSchemaGetter tableSchemaGetter,
	tableDeleter tableDeleter,
	tableFieldAdder tableFieldAdder,
	tableFieldDeleter tableFieldDeleter,
) tableHandler {
	return tableHandler{
		tableCreator,
		tableSchemaGetter,
		tableDeleter,
		tableFieldAdder,
		tableFieldDeleter,
	}
}

func (t *tableHandler) PutTable(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid table name"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var tableSchemaDto dto.TableSchemaDto
	json.Unmarshal(bodyBytes, &tableSchemaDto)

	tableSchemaResult := tableSchemaDto.ToModel()

	if tableSchemaResult.IsErr() {
		middleware.AttachError(w, tableSchemaResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid table schema"))
		return
	}

	err = t.tableCreator.CreateTable(
		projectId,
		tableNameResult.Unwrap(),
		tableSchemaResult.Unwrap(),
	)

	if err != nil {
		middleware.AttachError(w, err)

		if err, ok := err.(errs.InvalidTableError); ok {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating table"))
		return
	}
}

func (t *tableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid table name"))
		return
	}

	tableSchemaResult := t.tableSchemaGetter.GetTableSchema(projectId, tableNameResult.Unwrap())

	if tableSchemaResult.IsErr() {
		err := tableSchemaResult.UnwrapErr()

		middleware.AttachError(w, err)

		if _, ok := err.(errs.TableNotFoundError); ok {
			w.WriteHeader(404)
			w.Write([]byte("table not found"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting table"))
		return
	}

	tableSchemaDto := dto.GetTableSchemaDto(tableSchemaResult.Unwrap())

	resBodyBytes, _ := json.Marshal(tableSchemaDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (t *tableHandler) GetTables(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	tableSchemasResult := t.tableSchemaGetter.GetTableSchemas(projectId)

	if tableSchemasResult.IsErr() {
		err := tableSchemasResult.UnwrapErr()
		middleware.AttachError(w, err)

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting table"))
		return
	}

	tableSchemasDto := dto.GetTableSchemasDto(tableSchemasResult.Unwrap())

	resBodyBytes, _ := json.Marshal(tableSchemasDto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}

func (t *tableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid table name"))
		return
	}

	err := t.tableDeleter.DeleteTable(projectId, tableNameResult.Unwrap())

	if err != nil {
		middleware.AttachError(w, err)
		w.WriteHeader(500)
		w.Write([]byte("unexpected error deleting table"))
		return
	}

	w.WriteHeader(204)
}

func (t *tableHandler) AddField(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid table name"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var fieldCreationRequestDto dto.FieldCreationRequestDto
	json.Unmarshal(bodyBytes, &fieldCreationRequestDto)

	fieldCreationRequestResult := fieldCreationRequestDto.ToModel()

	if fieldCreationRequestResult.IsErr() {
		middleware.AttachError(w, fieldCreationRequestResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid request body"))
		return
	}

	fieldCreationRequest := fieldCreationRequestResult.Unwrap()

	err = t.tableFieldAdder.AddField(
		projectId,
		tableNameResult.Unwrap(),
		fieldCreationRequest.Name,
		fieldCreationRequest.Definition,
		fieldCreationRequest.DefaultValue,
	)

	if err != nil {
		middleware.AttachError(w, err)

		if err, ok := err.(errs.MissingDefaultValue); ok {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating field"))
		return
	}

	w.WriteHeader(200)
}

func (t *tableHandler) DeleteField(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid table name"))
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		panic("error reading body")
	}

	var fieldDeletionRequestDto dto.FieldDeletionRequestDto
	json.Unmarshal(bodyBytes, &fieldDeletionRequestDto)

	fieldDeletionRequestResult := fieldDeletionRequestDto.ToModel()

	if fieldDeletionRequestResult.IsErr() {
		middleware.AttachError(w, fieldDeletionRequestResult.UnwrapErr())
		w.WriteHeader(400)
		w.Write([]byte("invalid request body"))
		return
	}

	fieldDeletionRequest := fieldDeletionRequestResult.Unwrap()

	err = t.tableFieldDeleter.DeleteField(
		projectId,
		tableNameResult.Unwrap(),
		fieldDeletionRequest.Name,
	)

	if err != nil {
		middleware.AttachError(w, err)

		if _, ok := err.(errs.FieldNotFoundError); ok {
			w.WriteHeader(400)
			w.Write([]byte("field not found"))
			return
		}

		w.WriteHeader(500)
		w.Write([]byte("unexpected error deleting field"))
		return
	}

	w.WriteHeader(200)
}
