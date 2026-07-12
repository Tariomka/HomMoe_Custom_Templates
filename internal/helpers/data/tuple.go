package data

type Tuple[T1 any, T2 any] struct {
	Key   T1
	Value T2
}

func NewTuple[T1 any, T2 any](key T1, value T2) Tuple[T1, T2] {
	return Tuple[T1, T2]{Key: key, Value: value}
}
