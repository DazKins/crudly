package dto

import (
	"crudly/model"
	"crudly/util"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type EntityIdDto string

func (e EntityIdDto) ToModel() util.Result[model.EntityId] {
	uuid, err := uuid.Parse(string(e))

	if err != nil {
		return util.ResultErr[model.EntityId](errors.New("entity id is not a valid uuid"))
	}

	return util.ResultOk(model.EntityId(uuid))
}

type FieldDto any

func GetFieldDto(field model.Field) FieldDto {
	return FieldDto(any(field))
}

// Can't be ToModel since interface{} can't be extended
func FieldDtoToModel(f FieldDto) util.Result[model.Field] {
	return util.ResultOk(model.Field(any(f)))
}

type EntityDto map[string]FieldDto

func (e EntityDto) ToModel() util.Result[model.Entity] {
	result := model.Entity{}

	for k, v := range e {
		fieldResult := FieldDtoToModel(v)

		if fieldResult.IsErr() {
			err := fieldResult.UnwrapErr()
			return util.ResultErr[model.Entity](fmt.Errorf("error parsing field: %w", err))
		}

		result[k] = fieldResult.Unwrap()
	}

	return util.ResultOk(result)
}

func GetEntityDto(entity model.Entity) EntityDto {
	result := EntityDto{}

	for k, v := range entity {
		fieldDto := GetFieldDto(v)
		result[k] = fieldDto
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
