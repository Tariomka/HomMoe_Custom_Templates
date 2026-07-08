package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionIsNameless_AssignsManualName(t *testing.T) {
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	connection_editor.EnsureConnectionNames(connections)

	// Assert
	assert.Equal(t, "Manual-A-B", connections[0].Name)
}

func TestWhenConnectionAlreadyHasName_KeepsIt(t *testing.T) {
	// Arrange
	connections := []entities.Connection{
		{Name: "Rnd-A-B", From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	connection_editor.EnsureConnectionNames(connections)

	// Assert
	assert.Equal(t, "Rnd-A-B", connections[0].Name)
}

func TestWhenManualNameIsAlreadyTaken_AppendsNumericSuffix(t *testing.T) {
	// Arrange
	connections := []entities.Connection{
		{Name: "Manual-A-B", From: "Spawn-A", To: "Neutral-B"},
		{From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	connection_editor.EnsureConnectionNames(connections)

	// Assert
	assert.Equal(t, "Manual-A-B-2", connections[1].Name)
}

func TestWhenTwoNamelessConnectionsSharePair_AssignsDistinctNames(t *testing.T) {
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-B"},
		{From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	connection_editor.EnsureConnectionNames(connections)

	// Assert
	assert.NotEqual(t, connections[0].Name, connections[1].Name)
}
