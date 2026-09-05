package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneIsRemoved_KeepsOnlyOtherZones(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-C"},
		{Name: "Neutral-D"},
	}
	connections := []template_model.Connection{
		{From: "Spawn-A", To: "Neutral-C"},
	}

	// Act
	keptZones, _ := test_helpers.NewZoneEditorService().RemoveZone(zones, connections, "Neutral-C")

	// Assert
	assert.Equal(t, []template_model.Zone{{Name: "Spawn-A"}, {Name: "Neutral-D"}}, keptZones)
}

func TestWhenZoneIsRemoved_DropsConnectionsTouchingIt(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-C"},
		{Name: "Neutral-D"},
	}
	connections := []template_model.Connection{
		{From: "Spawn-A", To: "Neutral-C"},
		{From: "Neutral-C", To: "Neutral-D"},
		{From: "Spawn-A", To: "Neutral-D"},
	}

	// Act
	_, keptConnections := test_helpers.NewZoneEditorService().RemoveZone(zones, connections, "Neutral-C")

	// Assert
	assert.Equal(t, []template_model.Connection{{From: "Spawn-A", To: "Neutral-D"}}, keptConnections)
}
