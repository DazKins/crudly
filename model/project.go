package model

import "github.com/google/uuid"

type ProjectId uuid.UUID

func (p ProjectId) String() string {
	return uuid.UUID(p).String()
}

type ProjectKey string

type CreateProjectResponse struct {
	Id  ProjectId
	Key ProjectKey
}

type ProjectAuthInfo struct {
	Salt       string
	SaltedHash string
}
