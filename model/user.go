package model

import (
	"crudly/util/optional"

	"github.com/google/uuid"
)

type UserId uuid.UUID

func (u UserId) String() string {
	return uuid.UUID(u).String()
}

type User struct {
	TwitterId optional.O[string]
	GoogleId  optional.O[string]
}
