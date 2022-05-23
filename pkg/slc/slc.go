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

func Filter[T any](slc []T, f func(T) bool) (fltd []T) {
	for _, e := range slc {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return
}

func Map[T any](slc []T, f func(T) T) []T {
	mapped := make([]T, len(slc))

	for i, e := range slc {
		mapped[i] = f(e)
	}

	return mapped
}

func Find[T any](slc []T, f func(T) bool) (value T, found bool) {
	for _, e := range slc {
		if f(e) {
			value = e
			found = true
			return
		}
	}

	return
}

func RemoveByIndex[T any](slc []T, idx int) []T {
	return append(slc[:idx], slc[idx+1:]...)
}
