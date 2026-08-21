package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNoDebounceIsArmed_NextStateIsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	next := state.GetNextState()

	// Assert
	assert.Nil(t, next)
}

func TestWhenNextStateWasAssigned_NextStateMatchesTheAssignedValue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	assigned := editor_state_dto.NewDefaultEditorStateDto()
	assigned.PlayerCount++

	// Act
	state.SetNextState(assigned)

	// Assert
	assert.Equal(t, &assigned, state.GetNextState())
}

func TestWhenReturnedNextStateIsMutated_StoredPendingStateIsUnaffected(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(editor_state_dto.NewDefaultEditorStateDto())
	next := state.GetNextState()
	require.NotNil(t, next)

	// Act
	next.PlayerCount++

	// Assert
	assert.NotEqual(t, next.PlayerCount, state.GetNextState().PlayerCount)
}

func TestWhenNextStateIsReset_NextStateIsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(editor_state_dto.NewDefaultEditorStateDto())

	// Act
	state.ResetNextState()

	// Assert
	assert.Nil(t, state.GetNextState())
}
