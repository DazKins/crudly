package handler

import (
	"crudly/config"
	"crudly/http/dto"
	"crudly/model"
	"crudly/util"
	"encoding/json"
	"net/http"
)

type projectCreator interface {
	CreateProject() util.Result[model.CreateProjectResponse]
}

type adminProjectHandler struct {
	config         config.Config
	projectCreator projectCreator
}

func NewAdminProjectHandler(config config.Config, projectCreator projectCreator) adminProjectHandler {
	return adminProjectHandler{
		config,
		projectCreator,
	}
}

func (a adminProjectHandler) PostProject(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-api-key") != a.config.ProjectCreationApiKey {
		w.WriteHeader(401)
		return
	}

	createProjectResult := a.projectCreator.CreateProject()

	if createProjectResult.IsErr() {
		w.WriteHeader(500)
		return
	}

	dto := dto.GetCreateProjectResponseDto(createProjectResult.Unwrap())

	resBodyBytes, _ := json.Marshal(dto)

	w.Write(resBodyBytes)
}
