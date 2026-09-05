package file_service

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IFileService interface {
	LoadSettingsFile(filePath string) (*editor_state_model.EditorState, error)
	SaveSettings(filePath string, editorState *editor_state_model.EditorState) (string, error)
	SaveTemplateWithPreview(
		directory string,
		template *template_model.Template,
		previewImage *image.RGBA) (string, error)
}
