package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type ZoneEditorNeutralZoneRequestDto struct {
	Label         string
	Quality       neutral_zone.Quality
	CastleCount   int
	GenerateRoads bool
	Tuning        models.GenerationTuning
}
