package map_test

import (
	"sort"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapHasPairs_ReturnsQueryOverAllValues(t *testing.T) {
	t.Parallel()
	// Arrange
	source := map[string]int{"alpha": 10, "beta": 20, "gamma": 30}

	// Act
	values := linq.FromMap(source).SelectValues().ToSlice()

	// Assert
	sort.Ints(values)
	assert.Equal(t, []int{10, 20, 30}, values)
}

func TestWhenMapIsEmpty_ValueQueryYieldsNothing(t *testing.T) {
	t.Parallel()
	// Arrange
	source := map[string]int{}

	// Act
	values := linq.FromMap(source).SelectValues().ToSlice()

	// Assert
	assert.Empty(t, values)
}
