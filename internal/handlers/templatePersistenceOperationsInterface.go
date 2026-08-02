package handlers

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type ITemplatePersistenceOperations interface {
	SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error)
}
