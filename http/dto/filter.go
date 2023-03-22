package dto

import (
	"crudly/model"
	"crudly/util/result"
	"net/url"
	"strings"
)

func GetEntityFilterFromQuery(query url.Values) result.R[model.EntityFilter] {
	filterQuery := query["filter"]

	entityFilter := model.EntityFilter{}

	for _, v := range filterQuery {
		fieldFilterTypeResult := getFieldFilterTypeFromQuery(v)

		if fieldFilterTypeResult.IsErr() {
			return result.Err[model.EntityFilter](fieldFilterTypeResult.UnwrapErr())
		}

		fieldFilterType := fieldFilterTypeResult.Unwrap()

		vals := strings.Split(v, fieldFilterType.String())

		if len(vals) != 2 {
			return result.Errf[model.EntityFilter](
				"invalid filter: %s", v,
			)
		}

		fieldNameDto := FieldNameDto(vals[0])
		fieldNameResult := fieldNameDto.ToModel()

		if fieldNameResult.IsErr() {
			err := fieldNameResult.UnwrapErr()
			return result.Errf[model.EntityFilter]("error parsing filter field name: %w", err)
		}

		comparator := vals[1]

		entityFilter[fieldNameResult.Unwrap()] = model.FieldFilter{
			Type:       fieldFilterType,
			Comparator: comparator,
		}
	}

	return result.Ok(entityFilter)
}

func getFieldFilterTypeFromQuery(filterTypeQuery string) result.R[model.FieldFilterType] {
	if strings.Contains(filterTypeQuery, ">=") {
		return result.Ok(model.FieldFilterTypeGreaterThanEq)
	}

	if strings.Contains(filterTypeQuery, "<=") {
		return result.Ok(model.FieldFilterTypeLessThanEq)
	}

	if strings.Contains(filterTypeQuery, ">") {
		return result.Ok(model.FieldFilterTypeGreaterThan)
	}

	if strings.Contains(filterTypeQuery, "<") {
		return result.Ok(model.FieldFilterTypeLessThan)
	}

	if strings.Contains(filterTypeQuery, "=") {
		return result.Ok(model.FieldFilterTypeEquals)
	}

	return result.Errf[model.FieldFilterType]("invalid filter type: %s", filterTypeQuery)
}
