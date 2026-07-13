package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoPlayerLabelsProvided_ReturnsNeutralLabelsOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedChainZoneLabels(nil, mediumPlans("C", "D"))

	// Assert
	assert.Equal(t, []string{"C", "D"}, ordered)
}

func TestWhenSinglePlayerProvided_DistributesNeutralsAcrossEdgeGaps(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedChainZoneLabels([]string{"A"}, mediumPlans("C"))

	// Assert
	assert.Equal(t, []string{"A", "C"}, ordered)
}

func TestWhenTwoPlayersAndOneNeutralProvided_PlacesNeutralBetweenPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedChainZoneLabels([]string{"A", "B"}, mediumPlans("C"))

	// Assert
	assert.Equal(t, []string{"A", "C", "B"}, ordered)
}

func TestWhenMultipleNeutralsProvided_KeepsChainEndsOnPlayers(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedChainZoneLabels([]string{"A", "B"}, mediumPlans("C", "D", "E"))

	// Assert
	assert.Equal(t, []string{"A", "C", "E", "D", "B"}, ordered)
}
