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

type FieldDefinition struct {
	Type       FieldType
	Values     optional.O[[]string]
	PrimaryKey bool
}

type FieldName string

func (f FieldName) String() string {
	return string(f)
}

type TableSchema map[FieldName]FieldDefinition

type TableName string

func (t TableName) String() string {
	return string(t)
}

type TableSchemas map[TableName]TableSchema
