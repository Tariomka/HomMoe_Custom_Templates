package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTemplateNameWasUpdated_GetTemplateNameReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	templateName := gofakeit.ProductName()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.TemplateName = templateName })

	// Act
	actual := state.GetTemplateName()

	// Assert
	assert.Equal(t, templateName, actual)
}
