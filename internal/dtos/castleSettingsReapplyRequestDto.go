package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type CastleSettingsReapplyRequestDto struct {
	Zones       []entities.Zone
	Changes     editor_state_model.CastleSettingChanges
	EditorState editor_state_dto.EditorStateDto
}
