package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameHasHubPrefix_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Hub-B"

	// Act
	isHub := zone_helpers.IsZoneNameHub(zoneName)

	// Assert
	assert.True(t, isHub)
}

func TestWhenNameHasSpawnPrefix_IsZoneNameHubReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Spawn-B"

	// Act
	isHub := zone_helpers.IsZoneNameHub(zoneName)

	// Assert
	assert.False(t, isHub)
}
