package helpers

func MapSlice[TSource, TResult any](source []TSource, convert func(TSource) TResult) []TResult {
	if source == nil {
		return nil
	}

	result := make([]TResult, len(source))
	for index, item := range source {
		result[index] = convert(item)
	}

	return result
}

func GetUniqueElements[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if !seen[element] {
			uniqueSlice = append(uniqueSlice, element)
			seen[element] = true
		}
	}
	return uniqueSlice
}
