package model

type FieldSchema uint

const (
	FieldSchemaId      FieldSchema = 0
	FieldSchemaInteger FieldSchema = 1
	FieldSchemaString  FieldSchema = 2
	FieldSchemaBoolean FieldSchema = 3
	FieldSchemaTime    FieldSchema = 4
)

type TableName string

func (t TableName) String() string {
	return string(t)
}

type TableSchema map[string]FieldSchema
