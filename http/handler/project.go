package handler

import (
	"crudly/config"
	"crudly/http/dto"
	"crudly/http/middleware"
	"crudly/model"
	"crudly/util/result"
	"encoding/json"
	"net/http"
)

type projectCreator interface {
	CreateProject() result.R[model.CreateProjectResponse]
}

type projectHandler struct {
	config         config.Config
	projectCreator projectCreator
}

func NewProjectHandler(config config.Config, projectCreator projectCreator) projectHandler {
	return projectHandler{
		config,
		projectCreator,
	}
}

func (p projectHandler) PostProject(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-api-key") != p.config.ProjectCreationApiKey {
		w.WriteHeader(401)
		w.Write([]byte("incorrect api key"))
		return
	}

	createProjectResult := p.projectCreator.CreateProject()

	if createProjectResult.IsErr() {
		middleware.AttachError(w, createProjectResult.UnwrapErr())
		w.WriteHeader(500)
		w.Write([]byte("unexpected error creating project"))
		return
	}

	dto := dto.GetCreateProjectResponseDto(createProjectResult.Unwrap())

	resBodyBytes, _ := json.Marshal(dto)

	w.Header().Set("content-type", "application/json")
	w.Write(resBodyBytes)
}
