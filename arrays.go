package utils

func ArrayKeysDelete[T comparable](haystack []T, search ...T) []T {
	result := make([]T, 0, len(haystack))

	for _, v := range haystack {
		if !ArrayContains(search, v) {
			result = append(result, v)
		}
	}

	return result
}

func ArrayValueIndex[T comparable](haystack []T, search T) (int, bool) {
	for k, v := range haystack {
		if v == search {
			return k, true
		}
	}

	return 0, false
}

func ArrayContains[T comparable](haystack []T, search T) bool {
	for _, elem := range haystack {
		if elem == search {
			return true
		}
	}

	return false
}

func ArrayEquals[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for _, elem := range a {
		if !ArrayContains(b, elem) {
			return false
		}
	}

	return true
}
