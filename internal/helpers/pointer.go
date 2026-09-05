package helpers

func ClonePointer[T any](source *T) *T {
	if source == nil {
		return nil
	}

	value := *source
	return &value
}

func MapPointer[TSource, TResult any](source *TSource, convert func(TSource) TResult) *TResult {
	if source == nil {
		return nil
	}

	value := convert(*source)
	return &value
}
