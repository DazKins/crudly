package util

func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := []K{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	c := map[K]V{}
	for k, v := range m {
		c[k] = v
	}
	return c
}
