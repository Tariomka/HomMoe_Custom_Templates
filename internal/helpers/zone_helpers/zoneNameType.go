package zone_helpers

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
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

func GetZoneConnectionGuardQuality(
	zoneA, zoneB string,
	zones []entities.Zone,
	playerNames []string) neutral_zone.Quality {
	aIsPlayer := zoneA != "" && slices.Contains(playerNames, zoneA)
	bIsPlayer := zoneB != "" && slices.Contains(playerNames, zoneB)
	if aIsPlayer && bIsPlayer {
		return neutral_zone.QualityUnknown
	}

	zoneAQuality := GetZoneGuardQuality(zoneA, zones, playerNames)
	zoneBQuality := GetZoneGuardQuality(zoneB, zones, playerNames)
	return max(zoneAQuality, zoneBQuality)
}

func GetZoneGuardQuality(zoneName string, zones []entities.Zone, playerNames []string) neutral_zone.Quality {
	if zoneName == "" {
		return neutral_zone.QualityLow
	}

	zone, ok := linq.FromSlice(zones).First(func(z entities.Zone) bool { return z.Name == zoneName })
	if !ok {
		return neutral_zone.QualityLow
	}

	if slices.Contains(playerNames, zone.Name) {
		return neutral_zone.QualityUnknown
	}

	switch GetZoneTypeFromName(zone.Name) {
	case preview.ZoneTypeHub:
		return neutral_zone.QualityHighest
	case preview.ZoneTypePlayer:
		return neutral_zone.QualityUnknown
	case preview.ZoneTypeNeutral:
		return neutral_zone.GetQualityFrom(zone)
	case preview.ZoneTypeUnknown:
		fallthrough
	default:
		return neutral_zone.QualityUnknown
	}
}
