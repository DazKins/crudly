package validation

import (
	"crudly/model"
	"fmt"
)

type tableSchemaValidator struct{}

func NewTableSchemaValidator() tableSchemaValidator {
	return tableSchemaValidator{}
}

func (t *tableSchemaValidator) ValidateTableSchema(schema model.TableSchema) error {
	for k, v := range schema {
		if v.Type == model.FieldTypeEnum {
			if v.Values.IsNone() {
				return fmt.Errorf("enum type \"%s\" definition must include a values array", k)
			}
		} else {
			if v.Values.IsSome() {
				return fmt.Errorf("non enum type definition \"%s\" has a values array", k)
			}
		}
	}

	return nil
}
