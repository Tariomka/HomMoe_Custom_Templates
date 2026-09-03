package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"

type ZoneEditorConnectionRequestDto struct {
	From            string
	To              string
	Zones           []template_model.Zone
	PlayerZoneNames map[string]bool
}
