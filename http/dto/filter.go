package dto

import (
	"crudly/model"
	"crudly/util/result"
	"net/url"
	"strings"
)

func GetEntityFilterFromQuery(query url.Values) result.Result[model.EntityFilter] {
	filterQuery := query["filter"]

	entityFilter := model.EntityFilter{}

	for _, v := range filterQuery {
		filterVals := strings.Split(v, "|")

		if len(filterVals) != 3 {
			return result.Errf[model.EntityFilter]("invalid entity filer: %s", v)
		}

		fieldName := filterVals[0]

		filterTypeResult := getFieldFilterTypeFromQuery(filterVals[1])

		if filterTypeResult.IsErr() {
			return result.Errf[model.EntityFilter](
				"error getting filter type %w",
				filterTypeResult.UnwrapErr().Error(),
			)
		}

		comparator := filterVals[2]

		fieldFilter := model.FieldFilter{
			Type:       filterTypeResult.Unwrap(),
			Comparator: comparator,
		}

		entityFilter[fieldName] = fieldFilter
	}

	return result.Ok(entityFilter)
}

func getFieldFilterTypeFromQuery(filterTypeQuery string) result.Result[model.FieldFilterType] {
	switch strings.ToLower(filterTypeQuery) {
	case "equal":
		return result.Ok(model.FieldFilterTypeEquals)
	case ">":
		return result.Ok(model.FieldFilterTypeGreaterThan)
	case ">=":
		return result.Ok(model.FieldFilterTypeGreaterThanEq)
	case "<":
		return result.Ok(model.FieldFilterTypeLessThan)
	case "<=":
		return result.Ok(model.FieldFilterTypeLessThanEq)
	}
	return result.Errf[model.FieldFilterType]("invalid field filter type: %s", filterTypeQuery)
}
