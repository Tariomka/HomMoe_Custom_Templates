package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
)

func TestWhenOnlyManualEditsDiffer_ReportsEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := editor_state_dto.NewDefaultEditorStateDto()
	right := left
	right.ManualZones = []editor_state.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenANonManualFieldDiffers_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := editor_state_dto.NewDefaultEditorStateDto()
	right := left
	right.TemplateName = "Different Name"

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}
