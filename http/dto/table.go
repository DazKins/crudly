package dto

import (
	"crudly/model"
	"crudly/util/optional"
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
	case "enum":
		return result.Ok(model.FieldTypeEnum)
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
	case model.FieldTypeEnum:
		return FieldTypeDto("enum")
	}
	panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldType))
}

type FieldDefinitionDto struct {
	Type   FieldTypeDto `json:"type"`
	Values *[]string    `json:"values,omitempty"`
}

func (d FieldDefinitionDto) ToModel() result.Result[model.FieldDefinition] {
	fieldTypeResult := d.Type.ToModel()

	if fieldTypeResult.IsErr() {
		result.Err[model.FieldDefinition](fmt.Errorf("error parsing field type: %w", fieldTypeResult.UnwrapErr()))
	}

	fieldType := fieldTypeResult.Unwrap()

	return result.Ok(model.FieldDefinition{
		Type:   fieldType,
		Values: optional.FromPointer(d.Values),
	})
}

func GetFieldDefinitionDto(d model.FieldDefinition) FieldDefinitionDto {
	return FieldDefinitionDto{
		Type:   GetFieldTypeDto(d.Type),
		Values: optional.ToPointer(d.Values),
	}
}

type FieldNameDto string

func (f FieldNameDto) ToModel() result.Result[model.FieldName] {
	return result.Ok(model.FieldName(string(f)))
}

func GetFieldNameDto(f model.FieldName) FieldNameDto {
	return FieldNameDto(string(f))
}

type TableSchemaDto map[FieldNameDto]FieldDefinitionDto

func (e TableSchemaDto) ToModel() result.Result[model.TableSchema] {
	res := model.TableSchema{}

	for k, v := range e {
		fieldNameResult := k.ToModel()

		if fieldNameResult.IsErr() {
			err := fieldNameResult.UnwrapErr()
			return result.Err[model.TableSchema](fmt.Errorf("error with field name: %w", err))
		}

		fieldDefinitionResult := v.ToModel()

		if fieldDefinitionResult.IsErr() {
			err := fieldDefinitionResult.UnwrapErr()
			return result.Err[model.TableSchema](fmt.Errorf("error with field definition: %w", err))
		}

		res[fieldNameResult.Unwrap()] = fieldDefinitionResult.Unwrap()
	}

	return result.Ok(res)
}

func GetTableSchemaDto(schema model.TableSchema) TableSchemaDto {
	result := TableSchemaDto{}

	for k, v := range schema {
		fieldName := GetFieldNameDto(k)
		fieldDefinitionDto := GetFieldDefinitionDto(v)

		result[fieldName] = fieldDefinitionDto
	}

	return result
}

type TableNameDto string

func (t TableNameDto) ToModel() result.Result[model.TableName] {
	tableName := model.TableName(string(t))
	return result.Ok(tableName)
}
