package repositories

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type TemplateRepository struct{}

func NewTemplateRepository() IFileRepository[entities.RmgTemplate] {
	return &TemplateRepository{}
}

func (this *TemplateRepository) Load(filePath string) (entities.RmgTemplate, error) {
	return entities.RmgTemplate{}, common_errors.ErrNotImplemented
}

func (this *TemplateRepository) Save(directory string, filename string, entity entities.RmgTemplate) (string, error) {
	return "", nil
}
