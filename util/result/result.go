package result

import "fmt"

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		err:   nil,
	}
}

func Err[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}

func Errf[T any](str string, format ...any) Result[T] {
	return Result[T]{
		err: fmt.Errorf(str, format...),
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

func (r Result[T]) UnwrapOrDefault(def T) T {
	if r.IsErr() {
		return def
	}
	return r.value
}

func (r Result[T]) UnwrapErr() error {
	if r.IsOk() {
		panic("unwrap err called on ok result")
	}
	return r.err
}
