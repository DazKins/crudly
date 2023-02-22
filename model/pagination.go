package model

import "fmt"

type PaginationLimit uint

func (p PaginationLimit) String() string {
	return fmt.Sprintf("%+d", uint(p))
}

const DefaultPaginationLimit PaginationLimit = PaginationLimit(20)

type PaginationOffset uint

func (p PaginationOffset) String() string {
	return fmt.Sprintf("%+d", uint(p))
}

const DefaultPaginationOffset PaginationOffset = PaginationOffset(0)

type PaginationParams struct {
	Limit  PaginationLimit
	Offset PaginationOffset
}
