package repositories

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
)

type PreviewRepository struct{}

func NewPreviewRepository() IFileRepository[image.RGBA] {
	return &PreviewRepository{}
}

func (this *PreviewRepository) Load(filePath string) (image.RGBA, error) {
	return image.RGBA{}, common_errors.ErrNotImplemented
}

func (this *PreviewRepository) Save(directory string, filename string, entity image.RGBA) (string, error) {
	return "", nil
}
