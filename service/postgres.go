package service

import (
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func getPostgresTableName(projectId model.ProjectId, tableName model.TableName) string {
	return projectId.String() + "-table-" + tableName.String()
}

func getPostgresSchemaTableName(projectId model.ProjectId) string {
	return projectId.String() + "-tables"
}

func getPostgresDatatype(fieldType model.FieldType) string {
	switch fieldType {
	case model.FieldTypeId:
		return "uuid"
	case model.FieldTypeBoolean:
		return "boolean"
	case model.FieldTypeInteger:
		return "integer"
	case model.FieldTypeString:
		return "varchar"
	case model.FieldTypeTime:
		return "timestamp"
	case model.FieldTypeEnum:
		return "varchar"
	}
	panic(fmt.Sprintf("invalid field type has entered the system: %+v", fieldType))
}

const PostgresTimeFormat = "2006-01-02T15:04:05Z"

func getPostgresFieldValue(field any) result.R[string] {
	switch v := field.(type) {
	case uuid.UUID:
		return result.Ok("'" + v.String() + "'")
	case int:
		return result.Ok(fmt.Sprintf("'%d'", v))
	case string:
		return result.Ok("'" + v + "'")
	case bool:
		if v {
			return result.Ok("'true'")
		}
		return result.Ok("'false'")
	case time.Time:
		return result.Ok("'" + v.Format(PostgresTimeFormat) + "'")
	}
	return result.Err[string](fmt.Errorf("field: %+v has unsupported type", field))
}
