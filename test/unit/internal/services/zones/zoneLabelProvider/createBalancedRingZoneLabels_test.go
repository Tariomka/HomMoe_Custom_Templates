package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoPlayerLabelsProvided_ReturnsBalancedNeutralRing(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedRingZoneLabels(nil, mediumPlans("C", "D"), 0)

	// Assert
	assert.Equal(t, []string{"C", "D"}, ordered)
}

func TestWhenNoNeutralZonesProvided_ReturnsPlayerLabelsUnchanged(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedRingZoneLabels([]string{"A", "B"}, nil, 0)

	// Assert
	assert.Equal(t, []string{"A", "B"}, ordered)
}

func TestWhenPlayersAndNeutralsProvided_InterleavesNeutralsIntoGaps(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	ordered := provider.CreateBalancedRingZoneLabels([]string{"A", "B"}, mediumPlans("C", "D"), 0)

	// Assert
	assert.Equal(t, []string{"A", "C", "B", "D"}, ordered)
}
