package editorStateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenModelPointersAreMapped_PreservesNilAndValues(t *testing.T) {
	t.Parallel()
	// Arrange
	model := test_helpers.NewAllFieldsEditorStateModel()
	mapper := mappers.NewEditorStateMapper()

	// Act
	actual := []*string{nil, nil}
	if dto := mapper.ToDtoPointer(&model); dto != nil {
		actual[0] = &dto.TemplateName
	}
	if dto := mapper.ToDtoPointer(nil); dto != nil {
		actual[1] = &dto.TemplateName
	}

	// Assert
	assert.Equal(t, []*string{&model.TemplateName, nil}, actual)
}
