package pairSet_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant/misc"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPairIsAddedInAscendingOrder_StoresPairAsGiven(t *testing.T) {
	// Arrange
	set := misc.NewPairSet()

	// Act
	set.Add(1, 2)

	// Assert
	assert.True(t, (*set)[[2]int{1, 2}])
}

func TestWhenPairIsAddedInDescendingOrder_StoresNormalizedAscendingPair(t *testing.T) {
	// Arrange
	set := misc.NewPairSet()

	// Act
	set.Add(7, 3)

	// Assert
	assert.True(t, (*set)[[2]int{3, 7}])
}

func TestWhenSamePairIsAddedInBothOrders_KeepsSingleEntry(t *testing.T) {
	// Arrange
	set := misc.NewPairSet()
	first := gofakeit.Number(0, 100)
	second := first + gofakeit.Number(1, 100)

	// Act
	set.Add(first, second)
	set.Add(second, first)

	// Assert
	assert.Len(t, *set, 1)
}

func TestWhenDistinctPairsAreAdded_KeepsEntryPerPair(t *testing.T) {
	// Arrange
	set := misc.NewPairSet()

	// Act
	set.Add(0, 1)
	set.Add(1, 2)
	set.Add(0, 2)

	// Assert
	assert.Len(t, *set, 3)
}
