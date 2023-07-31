package utils

type IContains interface {
	string | int | int16 | int32 | int64 | int8 | uint | uint16 | uint32 | uint64 | uint8 | float64 | float32
}

func Contains[T IContains](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
