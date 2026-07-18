package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameHasNeutralPrefix_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Neutral-B"

	// Act
	isNeutral := zone_helpers.IsZoneNameNeutral(zoneName)

	// Assert
	assert.True(t, isNeutral)
}

func TestWhenNameHasSpawnPrefix_IsZoneNameNeutralReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Spawn-B"

	// Act
	isNeutral := zone_helpers.IsZoneNameNeutral(zoneName)

	// Assert
	assert.False(t, isNeutral)
}
