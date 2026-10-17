package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEntityIsMappedBackToAState_TheStateEqualsTheOriginal(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewEditorStateMapper()
	expected := test_helpers.NewAllFieldsEditorStateModel()

	// Act
	actual := mapper.ToModel(mapper.ToEntity(expected))

	// Assert
	assert.Equal(t, expected, actual)
}

// The loaded generator version decides whether manual edits still await
// conversion, so it is carried rather than overwritten.
func TestWhenAnEntityIsMappedBackToAState_TheLoadedGeneratorVersionIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewEditorStateMapper()
	entity := mapper.NewDefaultEntity()
	entity.GeneratorVersion = 7

	// Act
	state := mapper.ToModel(entity)

	// Assert
	assert.Equal(t, 7, state.GeneratorVersion)
}
