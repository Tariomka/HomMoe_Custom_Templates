package repositories

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
)

const templateExtension = ".rmg.json"

type TemplateRepository struct {
	writer *atomicFileWriter
}

func NewTemplateRepository() IFileRepository[template_entity.RmgTemplate] {
	return &TemplateRepository{writer: newAtomicFileWriter()}
}

func (this *TemplateRepository) Load(_ string, _ *template_entity.RmgTemplate) error {
	return common_errors.ErrNotImplemented
}

func (this *TemplateRepository) Save(directory, filename string, entity template_entity.RmgTemplate) (string, error) {
	return this.writer.WriteJSON(directory, filename, templateExtension, &entity)
}
