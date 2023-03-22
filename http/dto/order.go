package dto

import (
	"crudly/model"
	"crudly/util/result"
	"net/url"
	"strings"
)

func GetEntityOrderFromQuery(query url.Values) result.R[model.EntityOrder] {
	orderQueries := query["order"]

	entityOrder := model.EntityOrder{}

	for _, orderQuery := range orderQueries {
		split := strings.Split(orderQuery, "|")

		fieldName := model.FieldName(split[0])
		fieldOrderType := model.FieldOrderTypeAscending

		if len(split) == 2 {
			switch split[1] {
			case "asc":
				fieldOrderType = model.FieldOrderTypeAscending
			case "desc":
				fieldOrderType = model.FieldOrderTypeDescending
			default:
				return result.Errf[model.EntityOrder](
					"invalid order type for field \"%s\": \"%s\"",
					fieldName.String(),
					split[1],
				)
			}
		}

		entityOrder[fieldName] = fieldOrderType
	}

	return result.Ok(entityOrder)
}
