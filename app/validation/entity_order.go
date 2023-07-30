package validation

import (
	"crudly/model"
	"crudly/util"
	"fmt"
)

type entityOrderValidator struct{}

func NewEntityOrderValidator() entityOrderValidator {
	return entityOrderValidator{}
}

func (e *entityOrderValidator) ValidateEntityOrders(
	entityOrders model.EntityOrders,
	tableSchema model.TableSchema,
) error {

	for _, entityOrder := range entityOrders {
		fieldName := entityOrder.FieldName
		fieldOrder := entityOrder.Type

		fieldDefinition, ok := tableSchema[fieldName]

		if !ok {
			return fmt.Errorf("field: \"%s\" does not exist", fieldName)
		}

		if !isValidFieldOrder(fieldDefinition.Type, fieldOrder) {
			return fmt.Errorf(
				"order: \"%s\" is not valid for field type \"%s\"",
				fieldOrder.String(),
				fieldDefinition.Type.String(),
			)
		}
	}

	return nil
}

func isValidFieldOrder(fieldType model.FieldType, fieldOrderType model.FieldOrderType) bool {
	validityMap := map[model.FieldType][]model.FieldOrderType{
		model.FieldTypeId: {},
		model.FieldTypeInteger: {
			model.FieldOrderTypeAscending,
			model.FieldOrderTypeDescending,
		},
		model.FieldTypeString: {
			model.FieldOrderTypeAscending,
			model.FieldOrderTypeDescending,
		},
		model.FieldTypeBoolean: {},
		model.FieldTypeTime: {
			model.FieldOrderTypeAscending,
			model.FieldOrderTypeDescending,
		},
		model.FieldTypeEnum: {},
	}

	validOrders := validityMap[fieldType]

	return util.Contains(validOrders, fieldOrderType)
}
