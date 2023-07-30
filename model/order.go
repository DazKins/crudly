package model

type FieldOrderType uint

const (
	FieldOrderTypeAscending  FieldOrderType = 0
	FieldOrderTypeDescending FieldOrderType = 1
)

func (e FieldOrderType) String() string {
	switch e {
	case FieldOrderTypeAscending:
		return "ascending"
	case FieldOrderTypeDescending:
		return "descending"
	}
	panic("invalid field order type has entered the system in stringify!")
}

type EntityOrder struct {
	Type      FieldOrderType
	FieldName FieldName
}

type EntityOrders []EntityOrder
