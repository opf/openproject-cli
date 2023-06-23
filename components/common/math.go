package common

func Max[T int|int8|int16|int32|int64](a, b T) T {
	if b > a {
		return b
	}

	return a
}
