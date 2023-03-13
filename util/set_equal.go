package util

func SetEqual[T comparable](t0 []T, t1 []T) bool {
	t0Map := createHashTable(t0)
	t1Map := createHashTable(t1)

	for k := range t0Map {
		_, ok := t1Map[k]

		if !ok {
			return false
		}
	}

	for k := range t1Map {
		_, ok := t0Map[k]

		if !ok {
			return false
		}
	}

	return true
}

func SetSubtract[T comparable](t0 []T, t1 []T) []T {
	t0Map := createHashTable(t0)

	for _, v := range t1 {
		delete(t0Map, v)
	}

	return Keys(t0Map)
}

// Does set subtract just on the keys, doesn't do anything with the values
// Values are taken from t0
func MapSubtract[T comparable, U any, V any](t0 map[T]U, t1 map[T]V) map[T]U {
	ret := map[T]U{}

	for k0, v0 := range t0 {
		if _, ok := t1[k0]; !ok {
			ret[k0] = v0
		}
	}

	return ret
}

func createHashTable[T comparable](t []T) map[T]struct{} {
	res := map[T]struct{}{}

	for _, v := range t {
		res[v] = struct{}{}
	}

	return res
}
