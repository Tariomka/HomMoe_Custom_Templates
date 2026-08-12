package base

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type ITopologyConnectionService interface {
	CreateRandomPortalConnections(
		playerLabels, orderedLabels []string,
		tuning models.GenerationTuning,
		maxCount int,
		neutralZones neutral_zone.Plans) []entities.Connection

	CreateMissingPlayerConnections(
		playerLabels []string,
		allZones []entities.Zone,
		connections []entities.Connection,
		tuning models.GenerationTuning) []entities.Connection

	CreateMissingConnections(
		playerLabels, allLabels []string,
		positions models.Positions,
		allZones []entities.Zone,
		connections []entities.Connection,
		tuning models.GenerationTuning,
		neutralZones neutral_zone.Plans) []entities.Connection

	GetBorderGuardValue(
		labelA, labelB string,
		playerLabels []string,
		neutralZones neutral_zone.Plans,
		tuning models.GenerationTuning) int
}
