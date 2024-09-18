package helper

func SliceFilter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func SliceFilterOne[T any](s []T, f func(T) bool) (T, int) {
	var (
		res T
		idx int
	)
	for i, v := range s {
		if f(v) {
			res = v
			idx = i
			break
		}
	}

	return res, idx
}

func SliceMap[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func SliceMapIdx[T, U any](s []T, f func(int, T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(i, v)
	}
	return r
}

func SliceFind[T any](s []T, f func(T) bool) (T, bool) {
	for _, v := range s {
		if f(v) {
			return v, true
		}
	}

	return Empty[T](), false
}

func Slice2Map[K comparable, V any, T any](s []T, f func(T) (K, V)) map[K]V {
	r := make(map[K]V, len(s))
	for _, v := range s {
		k, val := f(v)
		r[k] = val
	}
	return r
}

func SliceUnique[T comparable](s []T) []T {
	m := make(map[T]struct{})
	r := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := m[v]; !ok {
			r = append(r, v)
			m[v] = struct{}{}
		}
	}
	return r
}
