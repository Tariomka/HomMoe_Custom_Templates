package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type CastleSettingsReapplyRequestDto struct {
	Zones       []template_model.Zone
	Changes     editor_state_model.CastleSettingChanges
	EditorState editor_state_dto.EditorStateDto
}
