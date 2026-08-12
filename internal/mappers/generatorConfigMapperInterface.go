package mappers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type IGeneratorConfigMapper interface {
	FromEditorState(editorState dtos.EditorStateDto) *config.GeneratorConfig
}
