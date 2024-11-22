package utils

func Abs[T int16 | int32 | int64 | int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
