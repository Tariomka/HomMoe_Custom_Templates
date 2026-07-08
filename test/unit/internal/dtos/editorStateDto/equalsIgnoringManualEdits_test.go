package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestWhenStatesAreFullyIdentical_ReportsEqual(t *testing.T) {
	// Arrange
	left := dtos.NewDefaultEditorStateDto()
	right := left

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenOnlyManualEditFieldsDiffer_ReportsEqual(t *testing.T) {
	// Arrange
	left := dtos.NewDefaultEditorStateDto()
	right := left
	right.ManualZones = []editor_state_dto.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}
	right.ManualConnections = []editor_state_dto.ManualConnectionSave{
		{Connection: entities.Connection{Name: "A-B"}, IsUserAdded: true},
	}

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenNonManualFieldDiffers_ReportsNotEqual(t *testing.T) {
	// Arrange
	left := dtos.NewDefaultEditorStateDto()
	right := left
	right.TemplateName = "Different Name"

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}
