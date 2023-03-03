package dto

import (
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"strings"
)

type FieldTypeDto string

func (t FieldTypeDto) ToModel() result.Result[model.FieldType] {
	switch strings.ToLower(string(t)) {
	case "id":
		return result.Ok(model.FieldTypeId)
	case "integer":
		return result.Ok(model.FieldTypeInteger)
	case "string":
		return result.Ok(model.FieldTypeString)
	case "boolean":
		return result.Ok(model.FieldTypeBoolean)
	case "time":
		return result.Ok(model.FieldTypeTime)
	}
	return result.Err[model.FieldType](fmt.Errorf("unrecognised field type: %s", string(t)))
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

func (d FieldDefinitionDto) ToModel() result.Result[model.FieldDefinition] {
	fieldTypeResult := d.Type.ToModel()

	if fieldTypeResult.IsErr() {
		result.Err[model.FieldDefinition](fmt.Errorf("error parsing field type: %w", fieldTypeResult.UnwrapErr()))
	}

	fieldType := fieldTypeResult.Unwrap()

	return result.Ok(model.FieldDefinition{
		Type: fieldType,
	})
}

func GetFieldDefinitionDto(d model.FieldDefinition) FieldDefinitionDto {
	return FieldDefinitionDto{
		Type: GetFieldTypeDto(d.Type),
	}
}

type TableSchemaDto map[string]FieldDefinitionDto

func (e TableSchemaDto) ToModel() result.Result[model.TableSchema] {
	res := model.TableSchema{}

	for k, v := range e {
		fieldDefinitionResult := v.ToModel()

		if fieldDefinitionResult.IsErr() {
			err := fieldDefinitionResult.UnwrapErr()
			return result.Err[model.TableSchema](fmt.Errorf("error with field definition: %w", err))
		}

		res[k] = fieldDefinitionResult.Unwrap()
	}

	return result.Ok(res)
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

func (t TableNameDto) ToModel() result.Result[model.TableName] {
	tableName := model.TableName(string(t))
	return result.Ok(tableName)
}
