package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenADtoIsMappedToAModel_RoundTripsAndStampsCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	model := test_helpers.NewAllFieldsEditorStateModel()
	mapper := mappers.NewEditorStateMapper()
	model.SchemaVersion = 0

	// Act
	actual := mapper.ToModel(mapper.ToDto(model))
	model.SchemaVersion = editor_state.CurrentEditorStateSchemaVersion

	// Assert
	assert.Equal(t, model, actual)
}
