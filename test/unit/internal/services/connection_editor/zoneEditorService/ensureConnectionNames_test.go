package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionIsNameless_AssignsManualName(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	test_helpers.NewZoneEditorService().EnsureConnectionNames(connections)

	// Assert
	assert.Equal(t, "Manual-A-B", connections[0].Name)
}

func TestWhenConnectionAlreadyHasName_KeepsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Rnd-A-B", From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	test_helpers.NewZoneEditorService().EnsureConnectionNames(connections)

	// Assert
	assert.Equal(t, "Rnd-A-B", connections[0].Name)
}

func TestWhenManualNameIsAlreadyTaken_AppendsNumericSuffix(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{Name: "Manual-A-B", From: "Spawn-A", To: "Neutral-B"},
		{From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	test_helpers.NewZoneEditorService().EnsureConnectionNames(connections)

	// Assert
	assert.Equal(t, "Manual-A-B-2", connections[1].Name)
}

func TestWhenTwoNamelessConnectionsSharePair_AssignsDistinctNames(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []template_model.Connection{
		{From: "Spawn-A", To: "Neutral-B"},
		{From: "Spawn-A", To: "Neutral-B"},
	}

	// Act
	test_helpers.NewZoneEditorService().EnsureConnectionNames(connections)

	// Assert
	assert.NotEqual(t, connections[0].Name, connections[1].Name)
}
