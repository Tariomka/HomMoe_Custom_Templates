package repositories

import (
	"image"
	"image/png"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
)

const previewExtension = ".png"

type PreviewRepository struct {
	writer *atomicFileWriter
}

func NewPreviewRepository() IFileRepository[image.RGBA] {
	return &PreviewRepository{writer: newAtomicFileWriter()}
}

func (this *PreviewRepository) Load(_ string, _ *image.RGBA) error {
	return common_errors.ErrNotImplemented
}

func (this *PreviewRepository) Save(directory string, filename string, entity image.RGBA) (string, error) {
	return this.writer.Write(directory, filename, previewExtension, func(file *os.File) error {
		return png.Encode(file, &entity)
	})
}
