package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenManualEditsWereStored_GetManualEditsReturnsThem(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	edits := newPopulatedManualEdits()
	state.SetManualEdits(edits)

	// Act
	stored := state.GetManualEdits()

	// Assert
	assert.Equal(t, edits, stored)
}

func TestWhenNoManualEditsWereStored_GetManualEditsReturnsEmptySettings(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	stored := state.GetManualEdits()

	// Assert
	assert.Equal(t, editor_state_model.ManualEditSettings{}, stored)
}

func TestWhenTheReturnedEditsAreMutated_TheStoredEditsAreUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	edits := newPopulatedManualEdits()
	state.SetManualEdits(edits)
	returned := state.GetManualEdits()

	// Act
	returned.ManualZones[0].FieldOverrides.GuardedContentPool[0] = "mutated"
	returned.ManualZones[0].Name = "mutated"

	// Assert
	assert.Equal(t, edits, state.GetManualEdits())
}
