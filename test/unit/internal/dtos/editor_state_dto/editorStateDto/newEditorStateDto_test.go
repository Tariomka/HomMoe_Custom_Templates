package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenAModelIsWrappedForPersistence_TheDtoCarriesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := test_helpers.NewAllFieldsEditorStateModel()

	// Act
	dto := editor_state_dto.NewEditorStateDto(state)

	// Assert
	assert.Equal(t, state, dto.EditorState)
}
