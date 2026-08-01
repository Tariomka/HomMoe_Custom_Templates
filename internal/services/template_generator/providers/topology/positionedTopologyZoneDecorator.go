package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type PositionedTopologyZoneDecorator func(
	zones []entities.Zone,
	allLabels []string,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
)
