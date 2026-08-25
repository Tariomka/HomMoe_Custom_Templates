package state_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTemplateNameWasUpdated_GetTemplateNameReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		test_helpers.NewFileSystemHandler(),
		test_helpers.NewRegenerationHandler(),
		newEditorStateMapper(),
		false)
	templateName := gofakeit.ProductName()
	state.UpdateState(func(dto *editor_state_model.EditorState) { dto.TemplateName = templateName })

	// Act
	actual := state.GetTemplateName()

	// Assert
	assert.Equal(t, templateName, actual)
}
