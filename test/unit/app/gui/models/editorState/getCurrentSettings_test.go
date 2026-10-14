package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenManualEditsAreStored_CurrentSettingsHaveNone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()

	// Act
	settings := state.GetCurrentSettings()

	// Assert
	assert.False(t, settings.HasManualEdits())
}

func TestWhenCurrentSettingsAreRead_TheyMatchTheCurrentStateWithoutManualEdits(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()
	expected := state.GetCurrentState()
	expected.ManualEditSettings = editor_state_model.ManualEditSettings{}

	// Act
	settings := state.GetCurrentSettings()

	// Assert
	assert.Equal(t, expected, settings)
}

func TestWhenCurrentSettingsAreRead_StoredManualEditsAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()

	// Act
	state.GetCurrentSettings()

	// Assert
	assert.True(t, state.HasManualEdits())
}

func TestWhenReturnedCurrentSettingsAreMutated_StoredStateIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()
	settings := state.GetCurrentSettings()

	// Act
	settings.PlayerZoneContentRows[0].Sid += "-changed"

	// Assert
	assert.NotEqual(t, settings.PlayerZoneContentRows[0].Sid, state.GetCurrentState().PlayerZoneContentRows[0].Sid)
}

// newEditorStateWithManualEdits returns a state holding one manual zone on top
// of the default settings, whose player content rows are never empty.
func newEditorStateWithManualEdits() *models.EditorState {
	state := newEditorState()
	state.SetManualEdits(
		editor_state_model.ManualEditSettingsFromSnapshot([]template_model.Zone{{Name: "Neutral-A"}}, nil),
	)
	return state
}
