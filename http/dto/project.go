package dto

import (
	"crudly/model"
	"crudly/util/result"
	"fmt"

	"github.com/google/uuid"
)

type ProjectIdDto string

func GetProjectIdDto(projectId model.ProjectId) ProjectIdDto {
	return ProjectIdDto(uuid.UUID(projectId).String())
}

func (p ProjectIdDto) ToModel() result.Result[model.ProjectId] {
	uuid, err := uuid.Parse(string(p))

	if err != nil {
		return result.Err[model.ProjectId](fmt.Errorf("invalid uuid: %w", err))
	}

	return result.Ok(model.ProjectId(uuid))
}

type ProjectKeyDto string

func GetProjectKeyDto(projectKey model.ProjectKey) ProjectKeyDto {
	return ProjectKeyDto(string(projectKey))
}

type CreateProjectResponseDto struct {
	Id         ProjectIdDto  `json:"id"`
	ProjectKey ProjectKeyDto `json:"projectKey"`
}

func GetCreateProjectResponseDto(createProjectResponse model.CreateProjectResponse) CreateProjectResponseDto {
	return CreateProjectResponseDto{
		Id:         GetProjectIdDto(createProjectResponse.Id),
		ProjectKey: GetProjectKeyDto(createProjectResponse.Key),
	}
}
