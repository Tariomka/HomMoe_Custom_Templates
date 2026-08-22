package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsCloned_CloneEqualsTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_dto.NewDefaultEditorStateDto()

	// Act
	clone := state.Clone()

	// Assert
	assert.Equal(t, state, clone)
}

// TestWhenAContentRowIsMutatedOnTheClone_SourceIsUnchanged proves the shim
// delegates to the model's deep clone rather than copying the shell.
func TestWhenAContentRowIsMutatedOnTheClone_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := editor_state_dto.NewDefaultEditorStateDto()
	state.HubZoneContentRows = []models.ZoneContentRowSave{{Sid: "row", Count: 1}}
	clone := state.Clone()

	// Act
	clone.HubZoneContentRows[0].Sid = "changed"

	// Assert
	assert.Equal(t, "row", state.HubZoneContentRows[0].Sid)
}
