package model

type FieldSchema uint

const (
	FieldSchemaInteger FieldSchema = 0
	FieldSchemaString  FieldSchema = 1
	FieldSchemaBoolean FieldSchema = 2
	FieldSchemaId      FieldSchema = 3
)

type TableName string

func (t TableName) String() string {
	return string(t)
}

type TableSchema map[string]FieldSchema
