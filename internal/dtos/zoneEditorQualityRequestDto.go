package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ZoneEditorQualityRequestDto struct {
	Zone        template_model.Zone
	Quality     neutral_zone.Quality
	CastleCount int
	Tuning      models.GenerationTuning
}
