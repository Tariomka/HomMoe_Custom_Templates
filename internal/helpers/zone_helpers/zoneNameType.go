package zone_helpers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

func GetZoneTypeFromName(zoneName string) preview.ZoneType {
	if strings.EqualFold(zoneName, constants.HubZoneName) ||
		strings.HasPrefix(zoneName, constants.HubZonePrefix) {
		return preview.ZoneTypeHub
	}

	if strings.HasPrefix(zoneName, constants.PlayerZonePrefix) {
		return preview.ZoneTypePlayer
	}

	if strings.HasPrefix(zoneName, constants.NeutralZonePrefix) {
		return preview.ZoneTypeNeutral
	}

	return preview.ZoneTypeUnknown
}

func IsZoneNameNeutral(zoneName string) bool {
	return GetZoneTypeFromName(zoneName) == preview.ZoneTypeNeutral
}

func IsZoneNamePlayer(zoneName string) bool {
	return GetZoneTypeFromName(zoneName) == preview.ZoneTypePlayer
}

func IsZoneNameHub(zoneName string) bool {
	return GetZoneTypeFromName(zoneName) == preview.ZoneTypeHub
}
