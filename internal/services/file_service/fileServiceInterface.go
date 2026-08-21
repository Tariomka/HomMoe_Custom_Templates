package file_service

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type IFileService interface {
	LoadSettingsFile(filePath string) (*editor_state_dto.EditorStateDto, error)
	SaveSettings(filePath string, editorState *editor_state_dto.EditorStateDto) (string, error)
	SaveTemplateWithPreview(
		directory string,
		template *entities.RmgTemplate,
		previewImage *image.RGBA) (string, error)
}
