package optional

type Optional[T any] struct {
	value   T
	present bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
	}
}

func None[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

func (o Optional[T]) IsNone() bool {
	return !o.IsSome()
}

func (o Optional[T]) IsSome() bool {
	return o.present
}

func (o Optional[T]) Unwrap() T {
	if !o.present {
		panic("optional.unwrap called on None")
	}
	return o.value
}
