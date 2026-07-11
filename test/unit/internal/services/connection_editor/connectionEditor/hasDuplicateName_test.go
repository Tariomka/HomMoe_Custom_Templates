package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnotherConnectionSharesName_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "main-road"},
		{From: "Neutral-1", To: "Neutral-2", Name: "main-road"},
	}

	// Act
	hasDuplicate := connection_editor.HasDuplicateName(connections, &connections[0])

	// Assert
	assert.True(t, hasDuplicate)
}

func TestWhenNamesDifferOnlyByCase_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "Main-Road"},
		{From: "Neutral-1", To: "Neutral-2", Name: "main-road"},
	}

	// Act
	hasDuplicate := connection_editor.HasDuplicateName(connections, &connections[1])

	// Assert
	assert.True(t, hasDuplicate)
}

func TestWhenNamesAreDistinct_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "alpha"},
		{From: "Neutral-1", To: "Neutral-2", Name: "beta"},
	}

	// Act
	hasDuplicate := connection_editor.HasDuplicateName(connections, &connections[0])

	// Assert
	assert.False(t, hasDuplicate)
}

func TestWhenCurrentConnectionIsNil_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: "alpha"},
	}

	// Act
	hasDuplicate := connection_editor.HasDuplicateName(connections, nil)

	// Assert
	assert.False(t, hasDuplicate)
}

func TestWhenCurrentNameIsEmpty_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{
		{From: "Spawn-A", To: "Neutral-1", Name: ""},
		{From: "Neutral-1", To: "Neutral-2", Name: ""},
	}

	// Act
	hasDuplicate := connection_editor.HasDuplicateName(connections, &connections[0])

	// Assert
	assert.False(t, hasDuplicate)
}
