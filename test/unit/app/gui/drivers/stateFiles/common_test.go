package stateFiles_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
)

func newEditorStateMapper() mappers.IEditorStateMapper {
	return mappers.NewEditorStateMapper()
}
