package linq

import "iter"

type Predicate[T any] = func(T) bool

// Iterable is an interface that has to be implemented by a custom collection to work with linq.
type Iterable[T any] interface {
	Iterate() iter.Seq[T]
}

// Query is the type returned from query functions. It can be iterated manually
// as shown in the example.
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

// func (this Query[T]) Select[TResult any](selector func(T) TResult) Query[TResult] { // With go1.26.3 this is not valid
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

func (this Query[T]) First(predicate Predicate[T]) (T, bool) {
	var result T
	found := false
	this.Iterate(func(item T) bool {
		if predicate(item) {
			result = item
			found = true
			return false
		}
		return true
	})
	return result, found
}

func (this Query[T]) Any() bool {
	found := false
	this.Iterate(func(item T) bool {
		found = true
		return false
	})
	return found
}
