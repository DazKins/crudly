package model

type FieldOrderType uint

const (
	OrderTypeAscending  FieldOrderType = 0
	OrderTypeDescending FieldOrderType = 1
)

type EntityOrder map[FieldName]FieldOrderType
