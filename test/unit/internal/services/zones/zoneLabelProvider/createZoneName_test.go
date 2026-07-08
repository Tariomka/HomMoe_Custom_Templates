package zoneLabelProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelBelongsToPlayer_PrefixesSpawn(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	zoneName := provider.CreateZoneName("B", []string{"A", "B"})

	// Assert
	assert.Equal(t, "Spawn-B", zoneName)
}

func TestWhenLabelIsNotAPlayerLabel_PrefixesNeutral(t *testing.T) {
	// Arrange
	provider := zones.NewZoneLabelProvider()

	// Act
	zoneName := provider.CreateZoneName("C", []string{"A", "B"})

	// Assert
	assert.Equal(t, "Neutral-C", zoneName)
}
