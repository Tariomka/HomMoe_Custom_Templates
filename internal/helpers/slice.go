package helpers

import "slices"

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

func Equal[Slice ~[]Entity, Entity comparable](left, right Slice) bool {
	return (left == nil) == (right == nil) && slices.Equal(left, right)
}

func EqualFunc[Slice1 ~[]Entity1, Slice2 ~[]Entity2, Entity1, Entity2 any](
	left Slice1,
	right Slice2,
	equalFunction func(Entity1, Entity2) bool) bool {
	return (left == nil) == (right == nil) && slices.EqualFunc(left, right, equalFunction)
}
