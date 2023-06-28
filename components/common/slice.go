package common

func Contains[T comparable](slice []T, value T) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}

	return false
}

func Reduce[T, M any](slice []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range slice {
		acc = f(acc, v)
	}
	return acc
}

func Filter[T any, M []T](slice M, f func(T) bool) M {
	return Reduce(
		slice,
		func(state M, value T) M {
			if f(value) {
				return append(state, value)
			} else {
				return state
			}
		},
		[]T{},
	)
}
