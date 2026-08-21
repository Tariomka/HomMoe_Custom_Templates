package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type ITemplateHandler interface {
	GenerateTemplate(stateDto editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error)
	UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error)
	ReapplyCastleSettings(request dtos.CastleSettingsReapplyRequestDto) []entities.Zone
	SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error)
}
