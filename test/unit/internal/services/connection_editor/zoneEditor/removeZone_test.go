package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsRemoved_KeepsOnlyOtherZones(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-C"},
		{Name: "Neutral-D"},
	}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-C"},
	}

	// Act
	keptZones, _ := connection_editor.RemoveZone(zones, connections, "Neutral-C")

	// Assert
	assert.Equal(t, []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-D"}}, keptZones)
}

func TestWhenZoneIsRemoved_DropsConnectionsTouchingIt(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-C"},
		{Name: "Neutral-D"},
	}
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-C"},
		{From: "Neutral-C", To: "Neutral-D"},
		{From: "Spawn-A", To: "Neutral-D"},
	}

	// Act
	_, keptConnections := connection_editor.RemoveZone(zones, connections, "Neutral-C")

	// Assert
	assert.Equal(t, []entities.Connection{{From: "Spawn-A", To: "Neutral-D"}}, keptConnections)
}
