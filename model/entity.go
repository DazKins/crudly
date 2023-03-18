package model

import "github.com/google/uuid"

type EntityId uuid.UUID

func (e EntityId) String() string {
	return uuid.UUID(e).String()
}

type Field any

type Entity map[FieldName]Field

type Entities []Entity

type PartialEntity map[FieldName]Field
