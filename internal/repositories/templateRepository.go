package repositories

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

const templateExtension = ".rmg.json"

type TemplateRepository struct {
	writer *atomicFileWriter
}

func NewTemplateRepository() IFileRepository[entities.RmgTemplate] {
	return &TemplateRepository{writer: newAtomicFileWriter()}
}

func (this *TemplateRepository) Load(_ string) (entities.RmgTemplate, error) {
	return entities.RmgTemplate{}, common_errors.ErrNotImplemented
}

func (this *TemplateRepository) Save(
	directory string,
	filename string,
	entity entities.RmgTemplate) (string, error) {
	return this.writer.WriteJSON(directory, filename, templateExtension, &entity)
}
