package zone_helpers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

func GetZoneTypeFromName(zoneName string) preview.ZoneType {
	if strings.EqualFold(zoneName, "Hub") || strings.HasPrefix(zoneName, "Hub-") {
		return preview.ZoneTypeHub
	}

	if strings.HasPrefix(zoneName, "Spawn-") {
		return preview.ZoneTypePlayer
	}

	if strings.HasPrefix(zoneName, "Neutral-") {
		return preview.ZoneTypeNeutralZone
	}

	return preview.ZoneTypeUnknown
}

func IsZoneNameNeutral(zoneName string) bool {
	return GetZoneTypeFromName(zoneName) == preview.ZoneTypeNeutralZone
}

func IsZoneNamePlayer(zoneName string) bool {
	return GetZoneTypeFromName(zoneName) == preview.ZoneTypePlayer
}

func IsZoneNameHub(zoneName string) bool {
	return GetZoneTypeFromName(zoneName) == preview.ZoneTypeHub
}
