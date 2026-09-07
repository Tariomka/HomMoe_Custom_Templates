package linq

import "iter"

type PredicateMap[TKey comparable, TValue any] = func(TKey, TValue) bool

// QueryMap is the type returned from query functions. It can be iterated manually
// using Iterate property. Example:
//
//	for key, value := range linq.FromMap(myMap).Iterate {
//		// use key and value
//	}
//
//	linq.FromMap(myMap).Iterate(func(key TKey, value TValue) bool {
//		// use key and value
//		return true // to continue iteration
//		return false // to stop iteration
//	})
type QueryMap[TKey comparable, TValue any] struct {
	Iterate iter.Seq2[TKey, TValue]
}

// FromMap initializes a linq query with a passed map.
func FromMap[M ~map[TKey]TValue, TKey comparable, TValue any](source M) QueryMap[TKey, TValue] {
	return QueryMap[TKey, TValue]{
		Iterate: func(yield func(TKey, TValue) bool) {
			for key, value := range source {
				if !yield(key, value) {
					break
				}
			}
		},
	}
}

func (this QueryMap[TKey, TValue]) Where(predicate PredicateMap[TKey, TValue]) QueryMap[TKey, TValue] {
	return QueryMap[TKey, TValue]{
		Iterate: func(yield PredicateMap[TKey, TValue]) {
			this.Iterate(func(key TKey, value TValue) bool {
				if predicate(key, value) {
					return yield(key, value)
				}
				return true
			})
		},
	}
}

func (this QueryMap[TKey, TValue]) SelectKeys() Query[TKey] {
	return Query[TKey]{
		Iterate: func(yield Predicate[TKey]) {
			this.Iterate(func(key TKey, _ TValue) bool {
				return yield(key)
			})
		},
	}
}

func (this QueryMap[TKey, TValue]) SelectValues() Query[TValue] {
	return Query[TValue]{
		Iterate: func(yield Predicate[TValue]) {
			this.Iterate(func(_ TKey, value TValue) bool {
				return yield(value)
			})
		},
	}
}

func (this QueryMap[TKey, TValue]) ToMap() map[TKey]TValue {
	result := make(map[TKey]TValue)

	this.Iterate(func(key TKey, value TValue) bool {
		result[key] = value
		return true
	})

	return result
}
