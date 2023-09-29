package dto

import (
	"crudly/model"
	"crudly/util/result"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EntityIdDto string

func (e EntityIdDto) ToModel() result.R[model.EntityId] {
	uuid, err := uuid.Parse(string(e))

	if err != nil {
		return result.Err[model.EntityId](errors.New("entity id is not a valid uuid"))
	}

	return result.Ok(model.EntityId(uuid))
}

type FieldDto any

const TimeFormat = "2006-01-02T15:04:05"

func GetFieldDto(field model.Field) FieldDto {
	if time, ok := field.(time.Time); ok {
		return time.Format(TimeFormat)
	}

	return FieldDto(any(field))
}

// Can't be ToModel since interface{} can't be extended
func FieldDtoToModel(f FieldDto) result.R[model.Field] {
	return result.Ok(model.Field(any(f)))
}

type EntityDto map[FieldNameDto]FieldDto

func (e EntityDto) ToModel() result.R[model.Entity] {
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

func (e EntitiesDto) ToModel() result.R[model.Entities] {
	entities := model.Entities{}

	for index, entityDto := range e {
		entityResult := entityDto.ToModel()

		if entityResult.IsErr() {
			err := entityResult.UnwrapErr()

			return result.Errf[model.Entities](
				"error parsing entity at index: %d, %w",
				index,
				err,
			)
		}

		entities = append(entities, entityResult.Unwrap())
	}

	return result.Ok(entities)
}

func GetEntitiesDto(entities model.Entities) EntitiesDto {
	result := EntitiesDto{}

	for _, entity := range entities {
		result = append(result, GetEntityDto(entity))
	}

	return result
}

type PartialEntityDto map[FieldNameDto]FieldDto

func (p PartialEntityDto) ToModel() result.R[model.PartialEntity] {
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

type GetEntitiesResponseDto struct {
	Entities   EntitiesDto `json:"entities"`
	TotalCount int         `json:"totalCount"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
}

func GetGetEntitiesResponseDto(entities model.GetEntitiesResponse) GetEntitiesResponseDto {
	return GetEntitiesResponseDto{
		Entities:   GetEntitiesDto(entities.Entities),
		TotalCount: int(entities.TotalCount),
		Limit:      int(entities.Limit),
		Offset:     int(entities.Offset),
	}
}
