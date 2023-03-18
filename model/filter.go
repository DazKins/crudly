package model

type FieldFilterType uint

const (
	FieldFilterTypeEquals        FieldFilterType = 0
	FieldFilterTypeGreaterThan   FieldFilterType = 1
	FieldFilterTypeGreaterThanEq FieldFilterType = 2
	FieldFilterTypeLessThan      FieldFilterType = 3
	FieldFilterTypeLessThanEq    FieldFilterType = 4
)

func (f FieldFilterType) String() string {
	switch f {
	case FieldFilterTypeEquals:
		return "="
	case FieldFilterTypeGreaterThan:
		return ">"
	case FieldFilterTypeGreaterThanEq:
		return ">="
	case FieldFilterTypeLessThan:
		return "<"
	case FieldFilterTypeLessThanEq:
		return "<="
	}
	panic("invalid field filter type has entered the system in stringify!")
}

type FieldFilter struct {
	Type       FieldFilterType
	Comparator interface{}
}

type EntityFilter map[FieldName]FieldFilter
