package manualConnectionSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenSaveListIsEmpty_ReturnsNilConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state_model.ManualConnectionSave

	// Act
	connections := editor_state_model.FromManualConnectionSaves(saves)

	// Assert
	assert.Nil(t, connections)
}

func TestWhenSavesCarryUserAddedFlags_RestoresEachFlagOntoConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state_model.ManualConnectionSave{
		{Connection: entities.Connection{Name: "A-B", From: "Zone A", To: "Zone B"}, IsUserAdded: true},
		{Connection: entities.Connection{Name: "B-C", From: "Zone B", To: "Zone C"}, IsUserAdded: false},
	}
	expected := []template_model.Connection{
		{Name: "A-B", From: "Zone A", To: "Zone B", IsUserAdded: true},
		{Name: "B-C", From: "Zone B", To: "Zone C", IsUserAdded: false},
	}

	// Act
	connections := editor_state_model.FromManualConnectionSaves(saves)

	// Assert
	assert.Equal(t, expected, connections)
}
