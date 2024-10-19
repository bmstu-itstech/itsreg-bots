package maps

func Join[K comparable, V any](a, b map[K]V) map[K]V {
	res := make(map[K]V, len(a)+len(b))
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

func Filter[K comparable, V any](a map[K]V, f func(K, V) bool) map[K]V {
	res := make(map[K]V, len(a))
	for k, v := range a {
		if f(k, v) {
			res[k] = v
		}
	}
	return res
}
