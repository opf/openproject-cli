package common

func Contains[T comparable](slice []T, value T) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	
	return false
}

func Reduce[T, M any](s []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}