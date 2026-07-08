package map_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSourceMapGiven_IterationYieldsAllPairs(t *testing.T) {
	// Arrange
	firstValue := gofakeit.Number(1, 1000)
	secondValue := gofakeit.Number(1, 1000)
	source := map[string]int{"first": firstValue, "second": secondValue}

	// Act
	collected := make(map[string]int)
	for key, value := range linq.FromMap(source).Iterate {
		collected[key] = value
	}

	// Assert
	assert.Equal(t, map[string]int{"first": firstValue, "second": secondValue}, collected)
}

func TestWhenYieldReturnsFalse_IterationStopsEarly(t *testing.T) {
	// Arrange
	source := map[string]int{"first": 1, "second": 2, "third": 3}

	// Act
	yieldCount := 0
	linq.FromMap(source).Iterate(func(key string, value int) bool {
		yieldCount++
		return false
	})

	// Assert
	assert.Equal(t, 1, yieldCount)
}

func TestWhenSourceMapIsEmpty_IterationYieldsNothing(t *testing.T) {
	// Arrange
	source := map[string]int{}

	// Act
	yieldCount := 0
	linq.FromMap(source).Iterate(func(key string, value int) bool {
		yieldCount++
		return true
	})

	// Assert
	assert.Equal(t, 0, yieldCount)
}
