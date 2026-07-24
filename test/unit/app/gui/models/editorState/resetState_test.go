package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWhenModifiedStateIsReset_DefaultValuesAreRestored(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) {
		dto.TemplateName = gofakeit.Name()
		dto.PlayerCount = gofakeit.Number(3, 8)
	})

	// Act
	state.ResetState()

	// Assert
	assert.Equal(t, dtos.NewDefaultEditorStateDto(), state.GetCurrentState())
}

func TestWhenStateWithSnapshotIsReset_PreviousStateIsDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SnapshotCurrentState()
	require.True(t, state.HasPreviousState())

	// Act
	state.ResetState()

	// Assert
	assert.False(t, state.HasPreviousState())
}

func TestWhenStateWithNextStateIsReset_NextStateIsDropped(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetNextState(state.GetCurrentState())
	require.True(t, state.HasNextState())

	// Act
	state.ResetState()

	// Assert
	assert.False(t, state.HasNextState())
}
