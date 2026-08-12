package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type ZoneEditorQualityRequestDto struct {
	Zone        entities.Zone
	Quality     neutral_zone.Quality
	CastleCount int
	Tuning      models.GenerationTuning
}
