package model

import "github.com/google/uuid"

type EntityId uuid.UUID

func (e EntityId) String() string {
	return uuid.UUID(e).String()
}

type Field any

type Entity map[string]Field

type Entities []Entity
