package model

type FieldFilterType uint

const (
	FieldFilterTypeEquals        FieldFilterType = 0
	FieldFilterTypeGreaterThan   FieldFilterType = 1
	FieldFilterTypeGreaterThanEq FieldFilterType = 2
	FieldFilterTypeLessThan      FieldFilterType = 3
	FieldFilterTypeLessThanEq    FieldFilterType = 4
)

type FieldFilter struct {
	Type       FieldFilterType
	Comparator interface{}
}

type EntityFilter map[string]FieldFilter
