package util

type Result[T any] struct {
	value T
	err   error
}

func ResultOk[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		err:   nil,
	}
}

func ResultErr[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return !r.IsOk()
}

func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic("unwrap called on error result")
	}
	return r.value
}

func (r Result[T]) UnwrapErr() error {
	if r.IsOk() {
		panic("unwrap err called on ok result")
	}
	return r.err
}
