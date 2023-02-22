package dto

import (
	"crudly/model"
	"crudly/util"
	"fmt"
	"strconv"
)

type PaginationLimitPathParam string

func (p PaginationLimitPathParam) ToModel() util.Result[model.PaginationLimit] {
	limit, err := strconv.Atoi(string(p))

	if err != nil {
		return util.ResultErr[model.PaginationLimit](fmt.Errorf("limit is not an integer"))
	}

	if limit < 0 {
		return util.ResultErr[model.PaginationLimit](fmt.Errorf("limit is less than 0"))
	}

	return util.ResultOk(model.PaginationLimit(uint(limit)))
}

type PaginationOffsetPathParam string

func (p PaginationOffsetPathParam) ToModel() util.Result[model.PaginationOffset] {
	offset, err := strconv.Atoi(string(p))

	if err != nil {
		return util.ResultErr[model.PaginationOffset](fmt.Errorf("offset is not an integer"))
	}

	if offset < 0 {
		return util.ResultErr[model.PaginationOffset](fmt.Errorf("offset is less than 0"))
	}

	return util.ResultOk(model.PaginationOffset(uint(offset)))
}
