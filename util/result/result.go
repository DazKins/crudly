package result

import "fmt"

type R[T any] struct {
	value T
	err   error
}

func Ok[T any](value T) R[T] {
	return R[T]{
		value: value,
		err:   nil,
	}
}

func Err[T any](err error) R[T] {
	return R[T]{
		err: err,
	}
}

func Errf[T any](str string, format ...any) R[T] {
	return R[T]{
		err: fmt.Errorf(str, format...),
	}
}

func (r R[T]) IsOk() bool {
	return r.err == nil
}

func (r R[T]) IsErr() bool {
	return !r.IsOk()
}

func (r R[T]) Unwrap() T {
	if r.IsErr() {
		panic("unwrap called on error result")
	}
	return r.value
}

func (r R[T]) UnwrapOrDefault(def T) T {
	if r.IsErr() {
		return def
	}
	return r.value
}

func (r R[T]) UnwrapErr() error {
	if r.IsOk() {
		panic("unwrap err called on ok result")
	}
	return r.err
}
