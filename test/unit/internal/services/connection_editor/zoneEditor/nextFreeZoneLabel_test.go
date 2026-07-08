package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneListIsEmpty_ReturnsLabelA(t *testing.T) {
	// Arrange

	// Act
	label := connection_editor.NextFreeZoneLabel(nil)

	// Assert
	assert.Equal(t, "A", label)
}

func TestWhenFirstLettersAreUsed_ReturnsNextFreeLetter(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-B"},
	}

	// Act
	label := connection_editor.NextFreeZoneLabel(zones)

	// Assert
	assert.Equal(t, "C", label)
}

func TestWhenSameLetterIsUsedAcrossPrefixes_CountsItOnce(t *testing.T) {
	// Arrange
	zones := []entities.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-A"},
	}

	// Act
	label := connection_editor.NextFreeZoneLabel(zones)

	// Assert
	assert.Equal(t, "B", label)
}
