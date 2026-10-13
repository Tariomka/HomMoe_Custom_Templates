package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingWasGeneratedYet_PreviousSettingsAreNil(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()

	// Act
	settings := state.GetPreviousSettings()

	// Assert
	assert.Nil(t, settings)
}

func TestWhenSnapshotHasManualEdits_PreviousSettingsMatchItWithoutThem(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()
	state.SnapshotCurrentState()
	expected := *state.GetPreviousState()
	expected.ManualEditSettings = editor_state_model.ManualEditSettings{}

	// Act
	settings := state.GetPreviousSettings()

	// Assert
	assert.Equal(t, &expected, settings)
}

func TestWhenReturnedPreviousSettingsAreMutated_StoredSnapshotIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()
	state.SnapshotCurrentState()
	settings := state.GetPreviousSettings()

	// Act
	settings.PlayerZoneContentRows[0].Sid += "-changed"

	// Assert
	assert.NotEqual(t, settings.PlayerZoneContentRows[0].Sid, state.GetPreviousState().PlayerZoneContentRows[0].Sid)
}
