package util

import (
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](t0 T, t1 T) T {
	if t0 < t1 {
		return t0
	} else {
		return t1
	}
}
