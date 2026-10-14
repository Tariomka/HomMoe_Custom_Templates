package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
)

type ITemplateHandler interface {
	GenerateTemplate(state editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error)
	UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error)

	// CaptureManualEdits records the edited layout against a fresh generation
	// from the given settings.
	CaptureManualEdits(request dtos.ManualEditCaptureRequestDto) (dtos.ManualEditCaptureDto, error)

	// ExpandManualEdits rebuilds the compact manual edits on the given freshly
	// generated template.
	ExpandManualEdits(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto

	// ConvertManualEdits records an all-overridden snapshot against a fresh
	// generation from its baseline settings.
	ConvertManualEdits(request dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error)

	SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error)
}
