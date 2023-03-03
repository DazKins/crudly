package dto

import (
	"crudly/model"
	"crudly/util"
	"fmt"
	"strings"
)

type FieldTypeDto string

func (t FieldTypeDto) ToModel() util.Result[model.FieldType] {
	switch strings.ToLower(string(t)) {
	case "id":
		return util.ResultOk(model.FieldTypeId)
	case "integer":
		return util.ResultOk(model.FieldTypeInteger)
	case "string":
		return util.ResultOk(model.FieldTypeString)
	case "boolean":
		return util.ResultOk(model.FieldTypeBoolean)
	case "time":
		return util.ResultOk(model.FieldTypeTime)
	}
	return util.ResultErr[model.FieldType](fmt.Errorf("unrecognised field type: %s", string(t)))
}

func GetFieldTypeDto(fieldType model.FieldType) FieldTypeDto {
	switch fieldType {
	case model.FieldTypeId:
		return FieldTypeDto("id")
	case model.FieldTypeBoolean:
		return FieldTypeDto("boolean")
	case model.FieldTypeInteger:
		return FieldTypeDto("integer")
	case model.FieldTypeString:
		return FieldTypeDto("string")
	case model.FieldTypeTime:
		return FieldTypeDto("time")
	}
	panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldType))
}

type FieldDefinitionDto struct {
	Type FieldTypeDto `json:"type"`
}

func (d FieldDefinitionDto) ToModel() util.Result[model.FieldDefinition] {
	fieldTypeResult := d.Type.ToModel()

	if fieldTypeResult.IsErr() {
		util.ResultErr[model.FieldDefinition](fmt.Errorf("error parsing field type: %w", fieldTypeResult.UnwrapErr()))
	}

	fieldType := fieldTypeResult.Unwrap()

	return util.ResultOk(model.FieldDefinition{
		Type: fieldType,
	})
}

func GetFieldDefinitionDto(d model.FieldDefinition) FieldDefinitionDto {
	return FieldDefinitionDto{
		Type: GetFieldTypeDto(d.Type),
	}
}

type TableSchemaDto map[string]FieldDefinitionDto

func (e TableSchemaDto) ToModel() util.Result[model.TableSchema] {
	result := model.TableSchema{}

	for k, v := range e {
		fieldDefinitionResult := v.ToModel()

		if fieldDefinitionResult.IsErr() {
			err := fieldDefinitionResult.UnwrapErr()
			return util.ResultErr[model.TableSchema](fmt.Errorf("error with field definition: %w", err))
		}

		result[k] = fieldDefinitionResult.Unwrap()
	}

	return util.ResultOk(result)
}

func GetTableSchemaDto(schema model.TableSchema) TableSchemaDto {
	result := TableSchemaDto{}

	for k, v := range schema {
		fieldDefinitionDto := GetFieldDefinitionDto(v)

		result[k] = fieldDefinitionDto
	}

	return result
}

type TableNameDto string

func (t TableNameDto) ToModel() util.Result[model.TableName] {
	tableName := model.TableName(string(t))
	return util.ResultOk(tableName)
}
