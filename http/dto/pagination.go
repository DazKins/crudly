package dto

import (
	"crudly/model"
	"crudly/util/result"
	"fmt"
	"strconv"
)

type PaginationLimitPathParam string

func (p PaginationLimitPathParam) ToModel() result.Result[model.PaginationLimit] {
	limit, err := strconv.Atoi(string(p))

	if err != nil {
		return result.Err[model.PaginationLimit](fmt.Errorf("limit is not an integer"))
	}

	if limit < 0 {
		return result.Err[model.PaginationLimit](fmt.Errorf("limit is less than 0"))
	}

	return result.Ok(model.PaginationLimit(uint(limit)))
}

type PaginationOffsetPathParam string

func (p PaginationOffsetPathParam) ToModel() result.Result[model.PaginationOffset] {
	offset, err := strconv.Atoi(string(p))

	if err != nil {
		return result.Err[model.PaginationOffset](fmt.Errorf("offset is not an integer"))
	}

	if offset < 0 {
		return result.Err[model.PaginationOffset](fmt.Errorf("offset is less than 0"))
	}

	return result.Ok(model.PaginationOffset(uint(offset)))
}
