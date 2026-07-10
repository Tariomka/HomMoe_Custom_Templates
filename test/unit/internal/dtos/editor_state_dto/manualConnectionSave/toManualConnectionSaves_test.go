package manualConnectionSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionListIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var connections []entities.Connection

	// Act
	saves := editor_state_dto.ToManualConnectionSaves(connections)

	// Assert
	assert.Nil(t, saves)
}

func TestWhenConnectionsCarryUserAddedFlags_PreservesEachFlagInSave(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{
		{Name: "A-B", From: "Zone A", To: "Zone B", IsUserAdded: true},
		{Name: "B-C", From: "Zone B", To: "Zone C", IsUserAdded: false},
	}
	expected := []editor_state_dto.ManualConnectionSave{
		{Connection: connections[0], IsUserAdded: true},
		{Connection: connections[1], IsUserAdded: false},
	}

	// Act
	saves := editor_state_dto.ToManualConnectionSaves(connections)

	// Assert
	assert.Equal(t, expected, saves)
}
