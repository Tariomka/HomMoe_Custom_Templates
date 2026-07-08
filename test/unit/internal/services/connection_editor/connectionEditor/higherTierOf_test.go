package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothZonesArePlayerZones_ReturnsPlayerToPlayer(t *testing.T) {
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	playerZoneNames := map[string]bool{"Spawn-A": true, "Spawn-B": true}

	// Act
	tier := connection_editor.HigherTierOf("Spawn-A", "Spawn-B", zones, playerZoneNames)

	// Assert
	assert.Equal(t, connection_editor.ZoneTierPlayerToPlayer, tier)
}

func TestWhenSecondZoneTierIsHigher_ReturnsSecondZoneTier(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-Gold", GuardedContentPool: []string{"pool_t4_x"}},
	}
	playerZoneNames := map[string]bool{"Spawn-A": true}

	// Act
	tier := connection_editor.HigherTierOf("Spawn-A", "Neutral-Gold", zones, playerZoneNames)

	// Assert
	assert.Equal(t, connection_editor.ZoneTierGold, tier)
}

func TestWhenFirstZoneTierIsHigher_ReturnsFirstZoneTier(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Neutral-Gold", GuardedContentPool: []string{"pool_t5_x"}},
		{Name: "Neutral-Bronze", GuardedContentPool: []string{"pool_t1_x"}},
	}
	playerZoneNames := map[string]bool{}

	// Act
	tier := connection_editor.HigherTierOf("Neutral-Gold", "Neutral-Bronze", zones, playerZoneNames)

	// Assert
	assert.Equal(t, connection_editor.ZoneTierGold, tier)
}
