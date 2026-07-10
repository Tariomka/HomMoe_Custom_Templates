package map_test

import (
	"sort"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapHasPairs_ReturnsQueryOverAllKeys(t *testing.T) {
	t.Parallel()
	// Arrange
	source := map[string]int{"alpha": 1, "beta": 2, "gamma": 3}

	// Act
	keys := linq.FromMap(source).SelectKeys().ToSlice()

	// Assert
	sort.Strings(keys)
	assert.Equal(t, []string{"alpha", "beta", "gamma"}, keys)
}

func TestWhenMapIsEmpty_KeyQueryYieldsNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	source := map[string]int{}

	// Act
	keys := linq.FromMap(source).SelectKeys().ToSlice()

	// Assert
	assert.Empty(t, keys)
}
