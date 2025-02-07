package validation

import (
	"crudly/model"
	"crudly/util"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type entityValidator struct{}

func NewEntityValidator() entityValidator {
	return entityValidator{}
}

func (e *entityValidator) ValidateEntity(entity model.Entity, tableSchema model.TableSchema) error {
	for k := range entity {
		fieldDefinition, ok := tableSchema[k]

		if !ok {
			return fmt.Errorf("field \"%s\" does not exist in table schema", k)
		}

		err := validateField(entity, k, fieldDefinition)

		if err != nil {
			return err
		}
	}

	missingFields := util.MapSubtract(tableSchema, entity)

	for fieldName, fieldDefinition := range missingFields {
		if fieldDefinition.IsOptional {
			continue
		}

		return fmt.Errorf("missing field: %s", fieldName.String())
	}

	return nil
}

const IncomingTimeFormat = "2006-01-02T15:04:05"

func validateField(
	entity model.Entity,
	fieldName model.FieldName,
	fieldDefinition model.FieldDefinition,
) error {
	field, ok := entity[fieldName]

	if !ok {
		return fmt.Errorf("entity is missing field: \"%s\"", fieldName)
	}

	if fieldDefinition.IsOptional && field == nil {
		return nil
	}

	switch fieldDefinition.Type {
	case model.FieldTypeId:
		stringVal, ok := field.(string)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid id", fieldName)
		}

		uuidVal, err := uuid.Parse(stringVal)

		if err != nil {
			return fmt.Errorf("error parsing field \"%s\" as a uuid: %w", fieldName, err)
		}

		entity[fieldName] = uuidVal

		return nil
	case model.FieldTypeInteger:
		floatVal, ok := field.(float64)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid integer", fieldName)
		}

		truncated := math.Trunc(floatVal)

		if truncated != floatVal {
			return fmt.Errorf("field: \"%s\" is not a valid integer", fieldName)
		}

		entity[fieldName] = int(truncated)

		return nil
	case model.FieldTypeBoolean:
		boolean, ok := field.(bool)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid boolean", fieldName)
		}

		entity[fieldName] = boolean

		return nil
	case model.FieldTypeString:
		_, ok := field.(string)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid string", fieldName)
		}

		return nil
	case model.FieldTypeTime:
		stringVal, ok := field.(string)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid time", fieldName)
		}

		timeResult := util.ValidateIncomingTime(stringVal)

		if timeResult.IsErr() {
			return fmt.Errorf("error parsing field \"%s\" as time: %w", fieldName, timeResult.UnwrapErr())
		}

		entity[fieldName] = timeResult.Unwrap()

		return nil
	case model.FieldTypeEnum:
		val, ok := field.(string)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid enum", fieldName)
		}

		values := fieldDefinition.Values.Unwrap()

		if !util.Contains(values, val) {
			return fmt.Errorf(
				"field: \"%s\" has value: \"%s\" which is not a supported value. supported values: %v",
				fieldName,
				val,
				values,
			)
		}

		return nil
	default:
		panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldDefinition.Type))
	}
}
