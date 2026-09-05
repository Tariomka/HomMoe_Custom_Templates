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

func IsZoneNameHub(zoneName string) bool { return GetZoneTypeFromName(zoneName) == preview.ZoneTypeHub }

// IsClusterHubZoneName reports whether the zone is one of the per-cluster hubs
// of a multi-hub topology ("Hub-A"). Unlike IsZoneNameHub it excludes the shared
// hub, so callers that count clusters do not count the shared hub as one.
func IsClusterHubZoneName(zoneName string) bool {
	return strings.HasPrefix(zoneName, constants.HubZonePrefix)
}

// IsSharedHubZoneName reports whether the zone is the single shared hub every
// other zone connects to ("Hub"). Unlike IsZoneNameHub it excludes the
// per-cluster hubs, so callers that look for one center never match several.
func IsSharedHubZoneName(zoneName string) bool {
	return zoneName == constants.HubZoneName
}
