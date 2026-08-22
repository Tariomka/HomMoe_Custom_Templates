package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultStateIsCreated_CarriesTheDefaultModel(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	state := editor_state_dto.NewDefaultEditorStateDto()

	// Assert
	assert.Equal(t, editor_state_model.NewDefaultEditorStateModel(), state.EditorStateModel)
}
