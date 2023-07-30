package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"crudly/ctx"
	"crudly/errs"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
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

type tableHandler struct {
	tableCreator      tableCreator
	tableSchemaGetter tableSchemaGetter
	tableDeleter      tableDeleter
	entityCountGetter entityCountGetter
}

type entityCountGetter interface {
	GetTotalEntityCount(
		projectId model.ProjectId,
		tableName model.TableName,
	) result.R[uint]
}

func NewTableHandler(
	tableCreator tableCreator,
	tableSchemaGetter tableSchemaGetter,
	tableDeleter tableDeleter,
	entityCountGetter entityCountGetter,
) tableHandler {
	return tableHandler{
		tableCreator,
		tableSchemaGetter,
		tableDeleter,
		entityCountGetter,
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

func (t *tableHandler) GetTotalEntityCount(w http.ResponseWriter, r *http.Request) {
	projectId := ctx.GetRequestProjectId(r)
	tableName := ctx.GetRequestTableName(r)

	totalEntityCountResult := t.entityCountGetter.GetTotalEntityCount(projectId, tableName)

	if totalEntityCountResult.IsErr() {
		err := totalEntityCountResult.UnwrapErr()
		middleware.AttachError(w, err)

		w.WriteHeader(500)
		w.Write([]byte("unexpected error getting total entity counts"))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"totalCount\":%d}", totalEntityCountResult.Unwrap())))
}
