package helpers

func ClonePointer[T any](source *T) *T {
	if source == nil {
		return nil
	}

	value := *source
	return &value
}
