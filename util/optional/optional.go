package optional

import (
	"encoding/json"
)

type O[T any] struct {
	value   T
	present bool
}

func Some[T any](value T) O[T] {
	return O[T]{
		value:   value,
		present: true,
	}
}

func None[T any]() O[T] {
	return O[T]{
		present: false,
	}
}

func (o O[T]) IsNone() bool {
	return !o.IsSome()
}

func (o O[T]) IsSome() bool {
	return o.present
}

func (o O[T]) Unwrap() T {
	if !o.present {
		panic("optional.unwrap called on None")
	}
	return o.value
}

func FromPointer[T any](p *T) O[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

func (o O[T]) ToPointer() *T {
	var res *T

	if o.IsSome() {
		res = new(T)
		*res = o.Unwrap()
	}

	return res
}

func (s O[T]) MarshalJSON() ([]byte, error) {
	if s.IsNone() {
		return json.Marshal(nil)
	} else {
		return json.Marshal(s.value)
	}
}

func (s *O[T]) UnmarshalJSON(bytes []byte) error {
	recv := new(T)

	err := json.Unmarshal(bytes, &recv)

	if err != nil {
		return err
	}

	if recv == nil {
		s.value = *new(T)
		s.present = false
	} else {
		s.value = *recv
		s.present = true
	}

	return nil
}
