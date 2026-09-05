package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IGeneratorConfigMapper interface {
	FromEditorState(editorState editor_state_model.EditorState) *config.GeneratorConfig
}
