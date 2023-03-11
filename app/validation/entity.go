package validation

import (
	"crudly/model"
	"crudly/util"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type entityValidator struct{}

func NewEntityValidator() entityValidator {
	return entityValidator{}
}

func (e entityValidator) ValidateEntity(entity model.Entity, tableSchema model.TableSchema) error {
	for k := range entity {
		err := validateField(entity, k, tableSchema[k])

		if err != nil {
			return err
		}
	}

	entityKeys := util.Keys(entity)
	tableSchemaKeys := util.Keys(tableSchema)

	keysEqual := util.SetEqual(entityKeys, tableSchemaKeys)

	if !keysEqual {
		return errors.New("table schema keys do not match up with entity keys")
	}

	return nil
}

const TimeFormat = "2006-01-02T15:04:05Z"

func validateField(
	entity model.Entity,
	fieldName string,
	fieldDefinition model.FieldDefinition,
) error {
	field, ok := entity[fieldName]

	if !ok {
		return fmt.Errorf("entity is missing field: \"%s\"", fieldName)
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

		time, err := time.Parse(TimeFormat, stringVal)

		if err != nil {
			return fmt.Errorf("error parsing field \"%s\" as time: %w", fieldName, err)
		}

		entity[fieldName] = time

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
