package dto

import (
	"crudly/model"
	"crudly/util"
	"fmt"
	"strings"
)

type FieldSchemaDto string

func (e FieldSchemaDto) ToModel() util.Result[model.FieldSchema] {
	switch strings.ToLower(string(e)) {
	case "integer":
		return util.ResultOk(model.FieldSchemaInteger)
	case "string":
		return util.ResultOk(model.FieldSchemaString)
	case "boolean":
		return util.ResultOk(model.FieldSchemaBoolean)
	}
	return util.ResultErr[model.FieldSchema](fmt.Errorf("unrecognised field schema: %s", string(e)))
}

func GetFieldSchemaDto(schema model.FieldSchema) FieldSchemaDto {
	switch schema {
	case model.FieldSchemaBoolean:
		return FieldSchemaDto("boolean")
	case model.FieldSchemaInteger:
		return FieldSchemaDto("integer")
	case model.FieldSchemaString:
		return FieldSchemaDto("string")
	case model.FieldSchemaId:
		return FieldSchemaDto("id")
	}
	panic(fmt.Sprintf("invalid field schema has entered the system: %+v", schema))
}

type TableSchemaDto map[string]FieldSchemaDto

func (e TableSchemaDto) ToModel() util.Result[model.TableSchema] {
	result := model.TableSchema{}

	for k, v := range e {
		fieldSchemaResult := v.ToModel()

		if fieldSchemaResult.IsErr() {
			err := fieldSchemaResult.UnwrapErr()
			return util.ResultErr[model.TableSchema](fmt.Errorf("error with field schema: %w", err))
		}

		result[k] = fieldSchemaResult.Unwrap()
	}

	return util.ResultOk(result)
}

func GetTableSchemaDto(schema model.TableSchema) TableSchemaDto {
	result := TableSchemaDto{}

	for k, v := range schema {
		fieldSchemaDto := GetFieldSchemaDto(v)

		result[k] = fieldSchemaDto
	}

	return result
}

type TableNameDto string

func (t TableNameDto) ToModel() util.Result[model.TableName] {
	tableName := model.TableName(string(t))
	return util.ResultOk(tableName)
}
