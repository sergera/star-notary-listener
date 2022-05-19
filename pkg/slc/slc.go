package slc

func ShallowEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Filter[T any](slc []T, f func(T) bool) []T {
	fltd := make([]T, 0)
	for _, e := range slc {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}
