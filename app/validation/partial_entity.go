package validation

import (
	"crudly/model"
	"crudly/util"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type partialEntityValidator struct{}

func NewPartialEntityValidator() partialEntityValidator {
	return partialEntityValidator{}
}

func (p *partialEntityValidator) ValidatePartialEntity(partialEntity model.PartialEntity, tableSchema model.TableSchema) error {
	for k := range partialEntity {
		fieldDefinition, ok := tableSchema[k]

		if !ok {
			return fmt.Errorf("field \"%s\" does not exist in table schema", k)
		}

		err := validatePartialField(partialEntity, k, fieldDefinition)

		if err != nil {
			return err
		}
	}
	return nil
}

func validatePartialField(
	partialEntity model.PartialEntity,
	fieldName model.FieldName,
	fieldDefinition model.FieldDefinition,
) error {
	field := partialEntity[fieldName]

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

		partialEntity[fieldName] = uuidVal

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

		partialEntity[fieldName] = int(truncated)

		return nil
	case model.FieldTypeBoolean:
		boolean, ok := field.(bool)

		if !ok {
			return fmt.Errorf("field: \"%s\" is not a valid boolean", fieldName)
		}

		partialEntity[fieldName] = boolean

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

		partialEntity[fieldName] = time

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
