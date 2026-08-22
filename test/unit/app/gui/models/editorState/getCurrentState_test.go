package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingWasModified_DefaultStateIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	current := state.GetCurrentState()

	// Assert
	assert.Equal(t, editor_state_model.NewDefaultEditorStateModel(), current)
}

func TestWhenReturnedStateIsMutated_StoredStateStaysUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	copyOfState := state.GetCurrentState()

	// Act
	copyOfState.TemplateName = gofakeit.Name()
	copyOfState.PlayerCount = gofakeit.Number(3, 8)

	// Assert
	assert.Equal(t, editor_state_model.NewDefaultEditorStateModel(), state.GetCurrentState())
}
