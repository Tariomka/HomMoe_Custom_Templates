package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNothingWasGeneratedYet_PreviousStateIsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	previous := state.GetPreviousState()

	// Assert
	assert.Nil(t, previous)
}

func TestWhenStateWasSnapshotted_PreviousStateMatchesTheSnapshot(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	snapshotted := state.GetCurrentState()

	// Act
	state.SnapshotCurrentState()

	// Assert
	assert.Equal(t, &snapshotted, state.GetPreviousState())
}

// The snapshot is handed out as a copy so a caller cannot corrupt the state
// that change detection compares against.
func TestWhenReturnedPreviousStateIsMutated_StoredSnapshotIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	previous := state.GetPreviousState()
	require.NotNil(t, previous)

	// Act
	previous.PlayerCount++

	// Assert
	assert.NotEqual(t, previous.PlayerCount, state.GetPreviousState().PlayerCount)
}

func TestWhenPreviousStateIsRequestedTwice_SeparateCopiesAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	first, second := state.GetPreviousState(), state.GetPreviousState()

	// Assert
	assert.NotSame(t, first, second)
}

func TestWhenStateIsOverridden_PreviousStateIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()

	// Act
	state.OverrideState(editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	assert.Nil(t, state.GetPreviousState())
}
