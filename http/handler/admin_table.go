package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"crudly/errs"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util"

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
	GetTableSchema(projectId model.ProjectId, name model.TableName) util.Result[model.TableSchema]
}

type adminTableHandler struct {
	tableCreator      tableCreator
	tableSchemaGetter tableSchemaGetter
}

func NewAdminTableHandler(tableCreator tableCreator, tableSchemaGetter tableSchemaGetter) adminTableHandler {
	return adminTableHandler{
		tableCreator,
		tableSchemaGetter,
	}
}

func (e adminTableHandler) PutTable(w http.ResponseWriter, r *http.Request) {
	projectId := util.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
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
		return
	}

	err = e.tableCreator.CreateTable(
		projectId,
		tableNameResult.Unwrap(),
		tableSchemaResult.Unwrap(),
	)

	if err != nil {
		middleware.AttachError(w, err)
		w.WriteHeader(500)
		return
	}
}

func (e adminTableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	projectId := util.GetRequestProjectId(r)

	vars := mux.Vars(r)

	tableNameDto := dto.TableNameDto(vars["tableName"])

	tableNameResult := tableNameDto.ToModel()

	if tableNameResult.IsErr() {
		middleware.AttachError(w, tableNameResult.UnwrapErr())
		w.WriteHeader(400)
		return
	}

	tableSchemaResult := e.tableSchemaGetter.GetTableSchema(projectId, tableNameResult.Unwrap())

	if tableSchemaResult.IsErr() {
		err := tableSchemaResult.UnwrapErr()

		middleware.AttachError(w, err)

		if errors.As(err, new(errs.TableNotFoundError)) {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(500)
		return
	}

	tableSchemaDto := dto.GetTableSchemaDto(tableSchemaResult.Unwrap())

	resBodyBytes, _ := json.Marshal(tableSchemaDto)

	w.Write(resBodyBytes)
}
