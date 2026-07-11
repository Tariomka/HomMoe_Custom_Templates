package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenEveryConnectionEndpointExists_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}, {Name: "Neutral-1"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	hasErrors := connection_editor.ComputeHasErrors(zones, connections)

	// Assert
	assert.False(t, hasErrors)
}

func TestWhenFromZoneIsMissing_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Neutral-1"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1"}}

	// Act
	hasErrors := connection_editor.ComputeHasErrors(zones, connections)

	// Assert
	assert.True(t, hasErrors)
}

func TestWhenToZoneIsMissing_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: "Spawn-A"}}
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-99"}}

	// Act
	hasErrors := connection_editor.ComputeHasErrors(zones, connections)

	// Assert
	assert.True(t, hasErrors)
}
