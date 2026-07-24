package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/entities"

type ZoneEditorConnectionRequestDto struct {
	From            string
	To              string
	Zones           []entities.Zone
	PlayerZoneNames map[string]bool
}
