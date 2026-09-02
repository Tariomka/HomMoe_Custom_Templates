package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type IZoneTierService interface {
	GetQuality(zone entities.Zone) neutral_zone.Quality

	GetGuardQuality(
		zoneName string,
		zones []entities.Zone,
		playerNames []string) neutral_zone.Quality

	GetConnectionGuardQuality(
		zoneA, zoneB string,
		zones []entities.Zone,
		playerNames []string) neutral_zone.Quality
}
