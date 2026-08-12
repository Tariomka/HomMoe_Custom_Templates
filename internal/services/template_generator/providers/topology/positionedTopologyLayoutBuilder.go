package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

// PositionedTopologyLayoutBuilder builds a topology-specific layout: the zone
// labels in placement order, exactly one position per label, and the
// label-index pairs that should be connected.
type PositionedTopologyLayoutBuilder func(playerLabels []string, neutralZones neutral_zone.Plans) (
	allLabels []string, positions models.Positions, pairs []models.ConnectionIndexes)
