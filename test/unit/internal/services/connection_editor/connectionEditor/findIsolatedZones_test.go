package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasNoConnections_ReturnsThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-1"},
		{Name: "Neutral-2"},
	}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
	}

	// Act
	isolated := connection_editor.FindIsolatedZones(zones, connections)

	// Assert
	assert.Equal(t, []string{"Neutral-2"}, isolated)
}

func TestWhenEveryZoneIsConnected_ReturnsNoZones(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-1"},
	}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
	}

	// Act
	isolated := connection_editor.FindIsolatedZones(zones, connections)

	// Assert
	assert.Empty(t, isolated)
}
