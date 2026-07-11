package map_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenQueryHasPairs_ReturnsMapWithAllOfThem(t *testing.T) {
	t.Parallel()
	// Arrange
	firstValue := gofakeit.Number(1, 1000)
	secondValue := gofakeit.Number(1, 1000)
	source := map[string]int{"first": firstValue, "second": secondValue}

	// Act
	result := linq.FromMap(source).ToMap()

	// Assert
	assert.Equal(t, map[string]int{"first": firstValue, "second": secondValue}, result)
}

func TestWhenResultMapIsMutated_SourceMapStaysUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	source := map[string]int{"first": 10, "second": 20}

	// Act
	result := linq.FromMap(source).ToMap()
	result["first"] = 999

	// Assert
	assert.Equal(t, map[string]int{"first": 10, "second": 20}, source)
}

func TestWhenQueryIsEmpty_ReturnsEmptyNonNilMap(t *testing.T) {
	t.Parallel()
	// Arrange
	source := map[string]int{}

	// Act
	result := linq.FromMap(source).ToMap()

	// Assert
	assert.Equal(t, map[string]int{}, result)
}
