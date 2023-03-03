package optional

import "encoding/json"

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

func FromPointer[T any](p *T) Optional[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

func ToPointer[T any](o Optional[T]) *T {
	var res *T

	if o.IsSome() {
		res = new(T)
		*res = o.Unwrap()
	}

	return res
}

func (s Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value   T
		Present bool
	}{
		Value:   s.value,
		Present: s.present,
	})
}

func (s *Optional[T]) UnmarshalJSON(bytes []byte) error {
	recv := new(struct {
		Value   T
		Present bool
	})

	err := json.Unmarshal(bytes, &recv)

	if err != nil {
		return err
	}

	s.value = recv.Value
	s.present = recv.Present

	return nil
}
