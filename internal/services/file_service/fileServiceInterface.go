package file_service

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type IFileService interface {
	LoadSettingsFile(filePath string) (*dtos.EditorStateDto, error)
	SaveSettings(filePath string, editorState *dtos.EditorStateDto) (string, error)
	SaveTemplateWithPreview(
		directory string,
		template *entities.RmgTemplate,
		previewImage *image.RGBA) (string, error)
}
