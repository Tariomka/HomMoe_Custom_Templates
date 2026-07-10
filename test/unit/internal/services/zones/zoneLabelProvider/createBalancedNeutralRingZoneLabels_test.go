package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoNeutralZonesProvided_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	labels := provider.CreateBalancedNeutralRingZoneLabels(nil, 2)

	// Assert
	assert.Empty(t, labels)
}

func TestWhenSingleNeutralZoneProvided_ReturnsItsLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	labels := provider.CreateBalancedNeutralRingZoneLabels(mediumPlans("C"), 2)

	// Assert
	assert.Equal(t, []string{"C"}, labels)
}

func TestWhenMultipleNeutralZonesProvided_ReturnsEveryLabelExactlyOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	labels := provider.CreateBalancedNeutralRingZoneLabels(mediumPlans("C", "D", "E"), 2)

	// Assert
	assert.ElementsMatch(t, []string{"C", "D", "E"}, labels)
}
