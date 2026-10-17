package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingIsPending_NextSettingsAreNil(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()

	// Act
	settings := state.GetNextSettings()

	// Assert
	assert.Nil(t, settings)
}

func TestWhenPendingStateHasManualEdits_NextSettingsMatchItWithoutThem(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()
	state.SetNextState(state.GetCurrentState())
	expected := *state.GetNextState()
	expected.ManualEditSettings = editor_state_model.ManualEditSettings{}

	// Act
	settings := state.GetNextSettings()

	// Assert
	assert.Equal(t, &expected, settings)
}

func TestWhenReturnedNextSettingsAreMutated_StoredPendingStateIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorStateWithManualEdits()
	state.SetNextState(state.GetCurrentState())
	settings := state.GetNextSettings()

	// Act
	settings.PlayerZoneContentRows[0].Sid += "-changed"

	// Assert
	assert.NotEqual(t, settings.PlayerZoneContentRows[0].Sid, state.GetNextState().PlayerZoneContentRows[0].Sid)
}
