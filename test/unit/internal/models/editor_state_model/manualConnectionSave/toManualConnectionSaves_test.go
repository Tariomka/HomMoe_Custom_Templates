package manualConnectionSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionListIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var connections []template_model.Connection

	// Act
	saves := editor_state_model.ToManualConnectionSaves(connections)

	// Assert
	assert.Nil(t, saves)
}

func TestWhenConnectionsCarryUserAddedFlags_PreservesEachFlagInSave(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "A-B", From: "Zone A", To: "Zone B", IsUserAdded: true},
		{Name: "B-C", From: "Zone B", To: "Zone C", IsUserAdded: false},
	}
	expected := []editor_state_model.ManualConnectionSave{
		{Connection: template_model.ToConnectionEntity(connections[0]), IsUserAdded: true},
		{Connection: template_model.ToConnectionEntity(connections[1]), IsUserAdded: false},
	}

	// Act
	saves := editor_state_model.ToManualConnectionSaves(connections)

	// Assert
	assert.Equal(t, expected, saves)
}
