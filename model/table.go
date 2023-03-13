package model

import (
	"crudly/util/optional"
)

type FieldType uint8

const (
	FieldTypeId      FieldType = 0
	FieldTypeInteger FieldType = 1
	FieldTypeString  FieldType = 2
	FieldTypeBoolean FieldType = 3
	FieldTypeTime    FieldType = 4
	FieldTypeEnum    FieldType = 5
)

func (f FieldType) String() string {
	switch f {
	case FieldTypeId:
		return "id"
	case FieldTypeInteger:
		return "integer"
	case FieldTypeString:
		return "string"
	case FieldTypeBoolean:
		return "boolean"
	case FieldTypeTime:
		return "time"
	case FieldTypeEnum:
		return "enum"
	}
	panic("invalid field type has entered the system in stringify!")
}

type TableName string

func (t TableName) String() string {
	return string(t)
}

type FieldDefinition struct {
	Type       FieldType
	Values     optional.Optional[[]string]
	PrimaryKey bool
}

type TableSchema map[string]FieldDefinition
