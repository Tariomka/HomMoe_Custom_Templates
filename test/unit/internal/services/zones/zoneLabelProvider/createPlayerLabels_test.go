package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenThreePlayersRequested_ReturnsFirstThreeAlphabetLetters(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	labels := provider.CreatePlayerLabels(3)

	// Assert
	assert.Equal(t, []string{"A", "B", "C"}, labels)
}

func TestWhenZeroPlayersRequested_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	labels := provider.CreatePlayerLabels(0)

	// Assert
	assert.Empty(t, labels)
}

func TestWhenArbitraryPlayerCountRequested_ReturnsThatManyLabels(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()
	playerCount := gofakeit.Number(1, 8)

	// Act
	labels := provider.CreatePlayerLabels(playerCount)

	// Assert
	assert.Len(t, labels, playerCount)
}
