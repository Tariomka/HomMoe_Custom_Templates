package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IZoneTierService interface {
	GetQuality(zone entities.Zone) neutral_zone.Quality

	ResolveQuality(zone template_model.Zone) neutral_zone.Quality

	GetGuardQuality(
		zoneName string,
		zones []template_model.Zone,
		playerNames []string) neutral_zone.Quality

	GetConnectionGuardQuality(
		zoneA, zoneB string,
		zones []template_model.Zone,
		playerNames []string) neutral_zone.Quality
}
