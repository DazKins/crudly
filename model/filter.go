package model

type FieldFilterType uint

const (
	FieldFilterTypeEquals FieldFilterType = 0
)

type FieldFilter struct {
	Type       FieldFilterType
	Comparator interface{}
}

type EntityFilter map[string]FieldFilter
