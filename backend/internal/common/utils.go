package common

func ToSet[T comparable](slice []T) map[T]struct{} {
	m := make(map[T]struct{}, len(slice))
	for _, elem := range slice {
		m[elem] = struct{}{}
	}
	return m
}

func ToSetBy[T comparable, E any](slice []E, f func(E) T) map[T]struct{} {
	m := make(map[T]struct{}, len(slice))
	for _, elem := range slice {
		m[f(elem)] = struct{}{}
	}
	return m
}
