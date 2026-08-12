package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type ZoneEditorOptionsDto struct {
	Topology      config.MapTopology
	Tuning        models.GenerationTuning
	GenerateRoads bool
}
