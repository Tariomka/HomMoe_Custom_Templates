package linq

import "iter"

type Predicate[T any] = func(T) bool

// Iterable is an interface that has to be implemented by a custom collection to work with linq.
type Iterable[T any] interface {
	Iterate() iter.Seq[T]
}

// Query is the type returned from query functions. It can be iterated manually
// using Iterate property. Example:
//
//	for value := range linq.FromSlice(mySlice).Iterate {
//		// use value
//	}
//
//	linq.FromSlice(mySlice).Iterate(func(value T) bool {
//		// use value
//		return true // to continue iteration
//		return false // to stop iteration
//	})
type Query[T any] struct {
	Iterate iter.Seq[T]
}

// FromSlice initializes a linq query with a passed slice.
func FromSlice[S ~[]T, T any](source S) Query[T] {
	return Query[T]{
		Iterate: func(yield Predicate[T]) {
			for _, item := range source {
				if !yield(item) {
					return
				}
			}
		},
	}
}

// Where filters a collection of values based on a predicate.
func (this Query[T]) Where(predicate Predicate[T]) Query[T] {
	return Query[T]{
		Iterate: func(yield Predicate[T]) {
			this.Iterate(func(item T) bool {
				if predicate(item) {
					return yield(item)
				}

				return true
			})
		},
	}
}

// SelectString projects each element of a collection into a new string form.
func (this Query[T]) SelectString(selector func(T) string) Query[string] {
	return Query[string]{
		Iterate: func(yield func(string) bool) {
			this.Iterate(func(item T) bool {
				return yield(selector(item))
			})
		},
	}
}

// func (this Query[T]) Select[TResult any](selector func(T) TResult) Query[TResult] { // Should be valid with 1.27
// 	return Query[TResult]{
// 		Iterate: func(yield func(TResult) bool) {
// 			this.Iterate(func(item T) bool {
// 				return yield(selector(item))
// 			})
// 		},
// 	}
// }

// ToSlice executes the query and returns the results as a slice. The returned slice is a copy, not a reference.
func (this Query[T]) ToSlice() []T {
	var result []T

	this.Iterate(func(item T) bool {
		result = append(result, item)
		return true
	})

	return result
}

func (this Query[T]) FirstOrDefault(predicate Predicate[T]) T {
	var result T

	this.Iterate(func(item T) bool {
		if predicate(item) {
			result = item
			return false
		}

		return true
	})

	return result
}

func (this Query[T]) First(predicate Predicate[T]) (result T, ok bool) {
	this.Iterate(func(item T) bool {
		if predicate(item) {
			result = item
			ok = true
			return false
		}

		return true
	})

	return result, ok
}

func (this Query[T]) Any() bool {
	found := false

	this.Iterate(func(T) bool { found = true; return false })

	return found
}

func (this Query[T]) AnyFunc(predicate Predicate[T]) bool {
	found := false

	this.Iterate(func(item T) bool {
		if predicate(item) {
			found = true
			return false
		}

		return true
	})

	return found
}

func (this Query[T]) AllFunc(predicate Predicate[T]) bool {
	if !this.Any() {
		return false
	}

	all := true
	this.Iterate(func(item T) bool {
		if !predicate(item) {
			all = false
			return false
		}

		return true
	})

	return all
}

// Distinct method returns distinct elements from a collection. The result is an
// unordered collection that contains no duplicate values.
func (this Query[T]) Distinct() Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[any]bool)

			this.Iterate(func(item T) bool {
				if _, seen := set[item]; !seen {
					set[item] = true
					return yield(item)
				}

				return true
			})
		},
	}
}

// DistinctBy method returns distinct elements from a collection. This method
// executes selector function for each element to determine a value to compare.
// The result is an unordered collection that contains no duplicate values.
func (this Query[T]) DistinctBy(selector func(T) any) Query[T] {
	return Query[T]{
		Iterate: func(yield func(T) bool) {
			set := make(map[any]bool)

			this.Iterate(func(item T) bool {
				key := selector(item)

				if _, seen := set[key]; !seen {
					set[key] = true
					return yield(item)
				}

				return true
			})
		},
	}
}
