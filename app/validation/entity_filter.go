package validation

import (
	"crudly/model"
	"crudly/util"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type entityFilterValidator struct{}

func NewEntityFilterValidator() entityFilterValidator {
	return entityFilterValidator{}
}

func (e entityFilterValidator) ValidateEntityFilter(
	entityFilter model.EntityFilter,
	tableSchema model.TableSchema,
) error {

	for k := range entityFilter {
		fieldDefinition, ok := tableSchema[k]

		if !ok {
			return fmt.Errorf("field: \"%s\" does not exist", k)
		}

		err := validateFieldFilter(k, entityFilter, fieldDefinition)

		if err != nil {
			return fmt.Errorf("error validating field: \"%s\": %w", k, err)
		}
	}

	return nil
}

func isValidFieldFilter(fieldType model.FieldType, fieldFilterType model.FieldFilterType) bool {
	validityMap := map[model.FieldType][]model.FieldFilterType{
		model.FieldTypeId: {
			model.FieldFilterTypeEquals,
		},
		model.FieldTypeInteger: {
			model.FieldFilterTypeEquals,
			model.FieldFilterTypeGreaterThan,
			model.FieldFilterTypeGreaterThanEq,
			model.FieldFilterTypeLessThan,
			model.FieldFilterTypeLessThanEq,
		},
		model.FieldTypeString: {
			model.FieldFilterTypeEquals,
		},
		model.FieldTypeBoolean: {
			model.FieldFilterTypeEquals,
		},
		model.FieldTypeTime: {
			model.FieldFilterTypeEquals,
			model.FieldFilterTypeGreaterThan,
			model.FieldFilterTypeGreaterThanEq,
			model.FieldFilterTypeLessThan,
			model.FieldFilterTypeLessThanEq,
		},
		model.FieldTypeEnum: {
			model.FieldFilterTypeEquals,
		},
	}

	validFilters := validityMap[fieldType]

	return util.Contains(validFilters, fieldFilterType)
}

func validateFieldFilter(
	fieldName model.FieldName,
	entityFilter model.EntityFilter,
	fieldDefinition model.FieldDefinition,
) error {
	parsedFieldFilter := entityFilter[fieldName]

	comparator, ok := parsedFieldFilter.Comparator.(string)

	if !ok {
		panic("comparator was not a string.")
	}

	if !isValidFieldFilter(fieldDefinition.Type, parsedFieldFilter.Type) {
		return fmt.Errorf(
			"filter: \"%s\" is not valid for field type \"%s\"",
			parsedFieldFilter.Type.String(),
			fieldDefinition.Type.String(),
		)
	}

	switch fieldDefinition.Type {
	case model.FieldTypeId:
		uuidVal, err := uuid.Parse(comparator)

		if err != nil {
			return fmt.Errorf("filter comparator is not an id: %s", comparator)
		}

		parsedFieldFilter.Comparator = uuidVal
	case model.FieldTypeInteger:
		intNum, err := strconv.Atoi(comparator)

		if err != nil {
			return fmt.Errorf("filter comparator is not an integer: %s", comparator)
		}

		parsedFieldFilter.Comparator = intNum
	case model.FieldTypeString:
	case model.FieldTypeBoolean:
		if comparator == "true" {
			parsedFieldFilter.Comparator = true
		} else if comparator == "false" {
			parsedFieldFilter.Comparator = false
		} else {
			return fmt.Errorf("filter comparator is not a boolean: %s", comparator)
		}
	case model.FieldTypeTime:
		time, err := time.Parse(TimeFormat, comparator)

		if err != nil {
			return fmt.Errorf("filter comparator is not a timestamp: %s", comparator)
		}

		parsedFieldFilter.Comparator = time
	case model.FieldTypeEnum:
		vals := fieldDefinition.Values

		if !util.Contains(vals.Unwrap(), comparator) {
			return fmt.Errorf("filter comparator is not a valid enum value: %s", comparator)
		}
	default:
		panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldDefinition.Type))
	}

	entityFilter[fieldName] = parsedFieldFilter

	return nil
}
