package dto

import (
	"crudly/model"
	"crudly/util/result"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type EntityIdDto string

func (e EntityIdDto) ToModel() result.Result[model.EntityId] {
	uuid, err := uuid.Parse(string(e))

	if err != nil {
		return result.Err[model.EntityId](errors.New("entity id is not a valid uuid"))
	}

	return result.Ok(model.EntityId(uuid))
}

type FieldDto any

func GetFieldDto(field model.Field) FieldDto {
	return FieldDto(any(field))
}

// Can't be ToModel since interface{} can't be extended
func FieldDtoToModel(f FieldDto) result.Result[model.Field] {
	return result.Ok(model.Field(any(f)))
}

type EntityDto map[FieldNameDto]FieldDto

func (e EntityDto) ToModel() result.Result[model.Entity] {
	res := model.Entity{}

	for k, v := range e {
		fieldNameResult := k.ToModel()

		if fieldNameResult.IsErr() {
			err := fieldNameResult.UnwrapErr()
			return result.Err[model.Entity](fmt.Errorf("error parsing field name: %w", err))
		}

		fieldResult := FieldDtoToModel(v)

		if fieldResult.IsErr() {
			err := fieldResult.UnwrapErr()
			return result.Err[model.Entity](fmt.Errorf("error parsing field: %w", err))
		}

		res[fieldNameResult.Unwrap()] = fieldResult.Unwrap()
	}

	return result.Ok(res)
}

func GetEntityDto(entity model.Entity) EntityDto {
	result := EntityDto{}

	for k, v := range entity {
		fieldNameDto := GetFieldNameDto(k)
		fieldDto := GetFieldDto(v)
		result[fieldNameDto] = fieldDto
	}

	return result
}

type EntitiesDto []EntityDto

func GetEntitiesDto(entities model.Entities) EntitiesDto {
	result := EntitiesDto{}

	for _, entity := range entities {
		result = append(result, GetEntityDto(entity))
	}

	return result
}

type PartialEntityDto map[FieldNameDto]FieldDto

func (p PartialEntityDto) ToModel() result.Result[model.PartialEntity] {
	res := model.PartialEntity{}

	for k, v := range p {
		fieldNameResult := k.ToModel()

		if fieldNameResult.IsErr() {
			err := fieldNameResult.UnwrapErr()
			return result.Err[model.PartialEntity](fmt.Errorf("error parsing field name: %w", err))
		}

		fieldResult := FieldDtoToModel(v)

		if fieldResult.IsErr() {
			err := fieldResult.UnwrapErr()
			return result.Err[model.PartialEntity](fmt.Errorf("error parsing field: %w", err))
		}

		res[fieldNameResult.Unwrap()] = fieldResult.Unwrap()
	}

	return result.Ok(res)
}
