package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameHasSpawnPrefix_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Spawn-A"

	// Act
	isPlayer := zone_helpers.IsZoneNamePlayer(zoneName)

	// Assert
	assert.True(t, isPlayer)
}

func TestWhenNameHasNeutralPrefix_IsZoneNamePlayerReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Neutral-A"

	// Act
	isPlayer := zone_helpers.IsZoneNamePlayer(zoneName)

	// Assert
	assert.False(t, isPlayer)
}
