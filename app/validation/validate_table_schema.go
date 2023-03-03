package validation

import (
	"crudly/model"
	"errors"
)

type tableSchemaValidator struct{}

func NewTableSchemaValidator() tableSchemaValidator {
	return tableSchemaValidator{}
}

func (t tableSchemaValidator) ValidateTableSchema(schema model.TableSchema) error {
	for _, v := range schema {
		if v.Type == model.FieldTypeEnum {
			if v.Values.IsNone() {
				return errors.New("enum types must include a values array")
			}
		}
	}

	return nil
}
