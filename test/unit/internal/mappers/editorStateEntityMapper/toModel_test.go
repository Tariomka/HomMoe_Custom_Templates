package editorStateEntityMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEntityIsMappedBackToAState_TheStateEqualsTheOriginal(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewEditorStateEntityMapper()
	expected := test_helpers.NewAllFieldsEditorStateModel()

	// Act
	actual := mapper.ToModel(mapper.ToEntity(expected))

	// Assert
	assert.Equal(t, expected, actual)
}
