package utils

func MapClone[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}

	return result
}

func MapKeys[M, V comparable](m map[M]V) []M {
	keys := make([]M, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}

	return keys
}
