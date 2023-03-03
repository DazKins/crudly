package model

type FieldType uint8

const (
	FieldTypeId      FieldType = 0
	FieldTypeInteger FieldType = 1
	FieldTypeString  FieldType = 2
	FieldTypeBoolean FieldType = 3
	FieldTypeTime    FieldType = 4
)

type TableName string

func (t TableName) String() string {
	return string(t)
}

type FieldDefinition struct {
	Type FieldType
}

type TableSchema map[string]FieldDefinition
