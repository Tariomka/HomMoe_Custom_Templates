package base

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ITopologyConnectionService interface {
	CreateRandomPortalConnections(
		playerLabels, orderedLabels []string,
		tuning models.GenerationTuning,
		maxCount int,
		neutralZones neutral_zone.Plans) []template_model.Connection

	CreateMissingPlayerConnections(
		playerLabels []string,
		allZones []template_model.Zone,
		connections []template_model.Connection,
		tuning models.GenerationTuning) []template_model.Connection

	CreateMissingConnections(
		playerLabels, allLabels []string,
		positions models.Positions,
		allZones []template_model.Zone,
		connections []template_model.Connection,
		tuning models.GenerationTuning,
		neutralZones neutral_zone.Plans) []template_model.Connection

	GetBorderGuardValue(
		labelA, labelB string,
		playerLabels []string,
		neutralZones neutral_zone.Plans,
		tuning models.GenerationTuning) int
}
