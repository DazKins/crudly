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

func createHashTable[T comparable](t []T) map[T]struct{} {
	res := map[T]struct{}{}

	for _, v := range t {
		res[v] = struct{}{}
	}

	return res
}
