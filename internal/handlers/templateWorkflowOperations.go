package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type TemplateWorkflowOperations interface {
	GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error)
	UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error)
	ReapplyCastleSettings(request dtos.CastleSettingsReapplyRequestDto) []entities.Zone
	ValidateEditorState(stateDto dtos.EditorStateDto, fixIssues bool) dtos.EditorStateValidationDto
}
