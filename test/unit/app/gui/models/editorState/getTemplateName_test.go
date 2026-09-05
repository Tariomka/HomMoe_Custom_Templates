package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTemplateNameWasUpdated_GetTemplateNameReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	templateName := gofakeit.ProductName()
	state.UpdateCurrentState(func(dto *editor_state_model.EditorState) { dto.TemplateName = templateName })

	// Act
	actual := state.GetTemplateName()

	// Assert
	assert.Equal(t, templateName, actual)
}
