package util

func Keys[T comparable, U any](m map[T]U) []T {
	res := []T{}

	for k := range m {
		res = append(res, k)
	}

	return res
}
