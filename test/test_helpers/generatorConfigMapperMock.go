package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/mock"
)

// GeneratorConfigMapperMock is a testify mock of
// mappers.IGeneratorConfigMapper, used to unit-test collaborators without the
// real editor-state-to-config mapping.
type GeneratorConfigMapperMock struct {
	mock.Mock
}

func (this *GeneratorConfigMapperMock) FromEditorState(
	editorState editor_state_model.EditorState,
) *config.GeneratorConfig {
	arguments := this.Called(editorState)
	configuration, _ := arguments.Get(0).(*config.GeneratorConfig)
	return configuration
}
