package common

func Contains[T comparable](slice []T, value T) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	
	return false
}
