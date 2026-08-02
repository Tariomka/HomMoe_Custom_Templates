package handlers

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type ITemplatePersistence interface {
	SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error)
}
