package zoneFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenSpawnZoneIsCreated_UsesPlayerLabelName(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewZoneFactory(nil, nil)

	// Act
	zone := factory.CreateSpawnZone("A", "Player1", nil, 0, false, 1, 0, true, newUnitTuning())

	// Assert
	assert.Equal(t, "Spawn-A", zone.Name)
}
