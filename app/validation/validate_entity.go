package validation

import (
	"crudly/model"
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
	for k, v := range tableSchema {
		err := validateField(entity, k, v)

		if err != nil {
			return err
		}
	}
	return nil
}

const TimeFormat = "2006-01-02 15:04:05"

func validateField(entity model.Entity, fieldName string, schema model.FieldSchema) error {
	field := entity[fieldName]

	switch schema {
	case model.FieldSchemaId:
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
	case model.FieldSchemaInteger:
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
	case model.FieldSchemaBoolean:
		boolean, ok := field.(bool)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid boolean", fieldName)
		}

		entity[fieldName] = boolean

		return nil
	case model.FieldSchemaString:
		_, ok := field.(string)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid id", fieldName)
		}

		return nil
	case model.FieldSchemaTime:
		stringVal, ok := field.(string)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid id", fieldName)
		}

		time, err := time.Parse(TimeFormat, stringVal)

		if err != nil {
			return fmt.Errorf("error parsing field \"%s\" as time: %w", fieldName, err)
		}

		entity[fieldName] = time

		return nil
	default:
		panic(fmt.Sprintf("invalid field schema has entered the system: %+v", schema))
	}
}
