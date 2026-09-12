package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type NeutralZoneQualityEditRequest struct {
	Zones           []template_model.Zone
	Connections     []template_model.Connection
	PlayerZoneNames []string

	ZoneName    string
	Quality     neutral_zone.Quality
	CastleCount int
	Tuning      GenerationTuning
}
