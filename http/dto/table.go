package dto

import (
	"crudly/model"
	"crudly/util/optional"
	"crudly/util/result"
	"fmt"
	"strings"
)

type FieldTypeDto string

func (t FieldTypeDto) ToModel() result.R[model.FieldType] {
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
	Type       FieldTypeDto `json:"type"`
	Values     *[]string    `json:"values,omitempty"`
	IsOptional bool         `json:"isOptional"`
}

func (d FieldDefinitionDto) ToModel() result.R[model.FieldDefinition] {
	fieldTypeResult := d.Type.ToModel()

	if fieldTypeResult.IsErr() {
		return result.Err[model.FieldDefinition](fmt.Errorf("error parsing field type: %w", fieldTypeResult.UnwrapErr()))
	}

	fieldType := fieldTypeResult.Unwrap()

	return result.Ok(model.FieldDefinition{
		Type:       fieldType,
		Values:     optional.FromPointer(d.Values),
		IsOptional: d.IsOptional,
	})
}

func GetFieldDefinitionDto(d model.FieldDefinition) FieldDefinitionDto {
	return FieldDefinitionDto{
		Type:       GetFieldTypeDto(d.Type),
		Values:     d.Values.ToPointer(),
		IsOptional: d.IsOptional,
	}
}

type FieldNameDto string

func (f FieldNameDto) ToModel() result.R[model.FieldName] {
	return result.Ok(model.FieldName(string(f)))
}

func GetFieldNameDto(f model.FieldName) FieldNameDto {
	return FieldNameDto(string(f))
}

type TableSchemaDto map[FieldNameDto]FieldDefinitionDto

func (e TableSchemaDto) ToModel() result.R[model.TableSchema] {
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

func (t TableNameDto) ToModel() result.R[model.TableName] {
	tableName := model.TableName(string(t))
	return result.Ok(tableName)
}

func GetTableNameDto(tableName model.TableName) TableNameDto {
	return TableNameDto(string(tableName))
}

type TableSchemasDto map[TableNameDto]TableSchemaDto

func GetTableSchemasDto(tableSchemas model.TableSchemas) TableSchemasDto {
	result := TableSchemasDto{}

	for tableName, tableSchema := range tableSchemas {
		result[GetTableNameDto(tableName)] = GetTableSchemaDto(tableSchema)
	}

	return result
}

type FieldCreationRequestDto struct {
	Name         FieldNameDto       `json:"name"`
	Definition   FieldDefinitionDto `json:"schema"`
	DefaultValue *any               `json:"defaultValue,omitempty"`
}

func (f FieldCreationRequestDto) ToModel() result.R[model.FieldCreationRequest] {
	nameResult := f.Name.ToModel()

	if nameResult.IsErr() {
		return result.Err[model.FieldCreationRequest](fmt.Errorf("error parsing field name: %w", nameResult.UnwrapErr()))
	}

	definitionResult := f.Definition.ToModel()

	if definitionResult.IsErr() {
		return result.Err[model.FieldCreationRequest](fmt.Errorf("error parsing field definition: %w", definitionResult.UnwrapErr()))
	}

	return result.Ok(model.FieldCreationRequest{
		Name:         nameResult.Unwrap(),
		Definition:   definitionResult.Unwrap(),
		DefaultValue: optional.FromPointer(f.DefaultValue),
	})
}

type FieldDeletionRequestDto struct {
	Name FieldNameDto `json:"name"`
}

func (f FieldDeletionRequestDto) ToModel() result.R[model.FieldDeletionRequest] {
	nameResult := f.Name.ToModel()

	if nameResult.IsErr() {
		return result.Err[model.FieldDeletionRequest](fmt.Errorf("error parsing field name: %w", nameResult.UnwrapErr()))
	}

	return result.Ok(model.FieldDeletionRequest{
		Name: nameResult.Unwrap(),
	})
}
