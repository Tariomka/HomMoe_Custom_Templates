package zone_helpers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

func GetZoneTypeFromName(zoneName string) preview.ZoneType {
	if strings.EqualFold(zoneName, common.HubZoneName) || strings.HasPrefix(zoneName, common.HubZonePrefix) {
		return preview.ZoneTypeHub
	}

	if strings.HasPrefix(zoneName, common.PlayerZonePrefix) {
		return preview.ZoneTypePlayer
	}

	if strings.HasPrefix(zoneName, common.NeutralZonePrefix) {
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
