package dto

import (
	"crudly/model"
	"crudly/util/result"
	"net/url"
	"strings"
)

func GetEntityOrderFromQuery(query url.Values) result.R[model.EntityOrders] {
	orderQueries := query["order"]

	entityOrders := model.EntityOrders{}

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
				return result.Errf[model.EntityOrders](
					"invalid order type for field \"%s\": \"%s\"",
					fieldName.String(),
					split[1],
				)
			}
		}

		entityOrders = append(entityOrders, model.EntityOrder{
			Type:      fieldOrderType,
			FieldName: fieldName,
		})
	}

	return result.Ok(entityOrders)
}
