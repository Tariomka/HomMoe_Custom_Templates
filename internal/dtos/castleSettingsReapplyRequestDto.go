package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type CastleSettingsReapplyRequestDto struct {
	Zones       []entities.Zone
	Changes     editor_state_dto.CastleSettingChanges
	EditorState editor_state_dto.EditorStateDto
}
