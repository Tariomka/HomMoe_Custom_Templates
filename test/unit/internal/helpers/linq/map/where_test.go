package map_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/stretchr/testify/assert"
)

func TestWhenPredicateMatchesSomePairs_KeepsOnlyMatchingOnes(t *testing.T) {
	// Arrange
	source := map[string]int{"first": 1, "second": 2, "third": 3, "fourth": 4}
	hasEvenValue := func(key string, value int) bool { return value%2 == 0 }

	// Act
	filtered := linq.FromMap(source).Where(hasEvenValue).ToMap()

	// Assert
	assert.Equal(t, map[string]int{"second": 2, "fourth": 4}, filtered)
}

func TestWhenPredicateMatchesNoPairs_ReturnsEmptyMap(t *testing.T) {
	// Arrange
	source := map[string]int{"first": 1, "third": 3}
	hasEvenValue := func(key string, value int) bool { return value%2 == 0 }

	// Act
	filtered := linq.FromMap(source).Where(hasEvenValue).ToMap()

	// Assert
	assert.Empty(t, filtered)
}

func TestWhenPredicateUsesKey_KeepsOnlyPairsWithMatchingKeys(t *testing.T) {
	// Arrange
	source := map[string]int{"keep-first": 1, "drop-second": 2, "keep-third": 3}
	keyStartsWithKeep := func(key string, value int) bool { return key[:4] == "keep" }

	// Act
	filtered := linq.FromMap(source).Where(keyStartsWithKeep).ToMap()

	// Assert
	assert.Equal(t, map[string]int{"keep-first": 1, "keep-third": 3}, filtered)
}
