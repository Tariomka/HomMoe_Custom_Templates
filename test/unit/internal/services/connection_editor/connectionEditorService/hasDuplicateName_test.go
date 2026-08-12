package connectionEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnotherConnectionSharesName_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "main-road"},
		{From: "Neutral-1", To: "Neutral-2", Name: "main-road"},
	}

	// Act
	hasDuplicate := service.HasDuplicateName(connections, &connections[0])

	// Assert
	assert.True(t, hasDuplicate)
}

func TestWhenNamesDifferOnlyByCase_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "Main-Road"},
		{From: "Neutral-1", To: "Neutral-2", Name: "main-road"},
	}

	// Act
	hasDuplicate := service.HasDuplicateName(connections, &connections[1])

	// Assert
	assert.True(t, hasDuplicate)
}

func TestWhenNamesAreDistinct_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "alpha"},
		{From: "Neutral-1", To: "Neutral-2", Name: "beta"},
	}

	// Act
	hasDuplicate := service.HasDuplicateName(connections, &connections[0])

	// Assert
	assert.False(t, hasDuplicate)
}

func TestWhenCurrentConnectionIsNil_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	connections := []entities.Connection{{From: "Spawn-A", To: "Neutral-1", Name: "alpha"}}

	// Act
	hasDuplicate := service.HasDuplicateName(connections, nil)

	// Assert
	assert.False(t, hasDuplicate)
}

func TestWhenCurrentNameIsEmpty_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneClassifier())
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1"},
		{From: "Neutral-1", To: "Neutral-2"},
	}

	// Act
	hasDuplicate := service.HasDuplicateName(connections, &connections[0])

	// Assert
	assert.False(t, hasDuplicate)
}
