package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNothingWasModified_DefaultStateIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	current := state.GetCurrentState()

	// Assert
	assert.Equal(t, dtos.NewDefaultEditorStateDto(), current)
}

func TestWhenReturnedStateIsMutated_StoredStateStaysUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	copyOfState := state.GetCurrentState()

	// Act
	copyOfState.TemplateName = gofakeit.Name()      //nolint:govet // this is testing object mutability
	copyOfState.PlayerCount = gofakeit.Number(3, 8) //nolint:govet // this is testing object mutability

	// Assert
	assert.Equal(t, dtos.NewDefaultEditorStateDto(), state.GetCurrentState())
}
