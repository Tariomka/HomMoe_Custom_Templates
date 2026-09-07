package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type PositionedTopologyZoneDecorator func(
	zones []template_model.Zone,
	allLabels []string,
	playerLabels []string,
	neutralZones neutral_zone.Plans)
