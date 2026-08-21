package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenStateIsOverridden_CurrentStateMatchesProvidedState(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	incoming := editor_state_dto.NewDefaultEditorStateDto()
	incoming.TemplateName = gofakeit.Name()
	incoming.PlayerCount = gofakeit.Number(3, 8)

	// Act
	state.OverrideState(incoming)

	// Assert
	assert.Equal(t, incoming, state.GetCurrentState())
}

func TestWhenStateWithSnapshotIsOverridden_PreviousStateIsDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	require.True(t, state.HasPreviousState())

	// Act
	state.OverrideState(editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	assert.False(t, state.HasPreviousState())
}

func TestWhenStateWithNextStateIsOverridden_NextStateIsDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())
	require.True(t, state.HasNextState())

	// Act
	state.OverrideState(editor_state_dto.NewDefaultEditorStateDto())

	// Assert
	assert.False(t, state.HasNextState())
}
