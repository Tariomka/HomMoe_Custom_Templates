package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IGeneratorConfigMapper interface {
	FromEditorState(editorState editor_state_dto.EditorStateDto) *config.GeneratorConfig
}
