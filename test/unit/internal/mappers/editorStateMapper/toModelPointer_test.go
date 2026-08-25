package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenDtoPointersAreMapped_PreservesNilAndValues(t *testing.T) {
	t.Parallel()
	// Arrange
	model := test_helpers.NewAllFieldsEditorStateModel()
	mapper := mappers.NewEditorStateMapper()
	dto := mapper.ToDto(model)

	// Act
	actual := []*string{nil, nil}
	if mappedModel := mapper.ToModelPointer(&dto); mappedModel != nil {
		actual[0] = &mappedModel.TemplateName
	}
	if mappedModel := mapper.ToModelPointer(nil); mappedModel != nil {
		actual[1] = &mappedModel.TemplateName
	}

	// Assert
	assert.Equal(t, []*string{&model.TemplateName, nil}, actual)
}
