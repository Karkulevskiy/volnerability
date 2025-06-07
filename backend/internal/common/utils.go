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

func Map[T1, T2 any](mapFn func(T1) T2, objs ...T1) []T2 {
	mapped := make([]T2, 0, len(objs))
	for _, obj := range objs {
		mapped = append(mapped, mapFn(obj))
	}
	return mapped
}

func Remove[T comparable](s []T, toRemove ...T) []T {
	clean := []T{}
	toRemoveSet := ToSet(toRemove)
	for _, v := range s {
		if _, ok := toRemoveSet[v]; !ok {
			clean = append(clean, v)
		}
	}
	return clean
}
