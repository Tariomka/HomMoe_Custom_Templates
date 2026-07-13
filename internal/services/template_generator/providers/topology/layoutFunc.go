package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
)

// layoutFunc builds a topology-specific layout: the zone labels in placement
// order (players first or interleaved, as the topology dictates), one position
// per label, and the label-index pairs that should be connected.
type layoutFunc func(playerLabels []string, neutralZones neutralZone.Plans) (
	allLabels []string, positions models.Positions, pairs []models.ConnectionIndexes)
