package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenEndpointTierIsGold_SeedsGoldGeneratorDefaults(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-Gold", GuardedContentPool: []string{"pool_t4_x"}},
	}
	playerZoneNames := map[string]bool{"Spawn-A": true}
	expected := entities.Connection{
		From:                 "Spawn-A",
		To:                   "Neutral-Gold",
		ConnectionType:       "Direct",
		GuardValue:           25000,
		GuardZone:            "Spawn-A",
		GuardMatchGroup:      "rnd_guard_A_Gold",
		GuardWeeklyIncrement: 0.15,
		IsUserAdded:          true,
	}

	// Act
	connection := connection_editor.NewDefaultConnection("Spawn-A", "Neutral-Gold", zones, playerZoneNames)

	// Assert
	assert.Equal(t, expected, connection)
}

func TestWhenBothEndpointsArePlayerZones_SeedsPlayerToPlayerGeneratorDefault(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Spawn-B"}}
	playerZoneNames := map[string]bool{"Spawn-A": true, "Spawn-B": true}

	// Act
	connection := connection_editor.NewDefaultConnection("Spawn-A", "Spawn-B", zones, playerZoneNames)

	// Assert
	assert.Equal(t, 30000, connection.GuardValue)
}
