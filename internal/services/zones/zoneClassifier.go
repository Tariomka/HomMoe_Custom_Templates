package zones

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type ZoneClassifier struct{}

func NewZoneClassifier() *ZoneClassifier {
	return &ZoneClassifier{}
}

func (this *ZoneClassifier) GetQuality(zone entities.Zone) neutral_zone.Quality {
	if strings.HasPrefix(zone.Name, constants.PlayerZonePrefix) {
		return neutral_zone.QualityUnknown
	}

	switch layoutValues := registry.GetLayoutValues(); zone.Layout {
	case layoutValues.Center:
		return this.getCenterQuality(zone)
	case layoutValues.TreasureZone:
		return this.getTreasureQuality(zone)
	case layoutValues.Sides:
		return this.getSidesQuality(zone)
	default:
		return neutral_zone.QualityUnknown
	}
}

func (this *ZoneClassifier) GetGuardQuality(
	zoneName string,
	zones []entities.Zone,
	playerNames []string,
) neutral_zone.Quality {
	if zoneName == "" {
		return neutral_zone.QualityLow
	}

	zone, ok := linq.FromSlice(zones).First(func(candidate entities.Zone) bool { return candidate.Name == zoneName })
	if !ok {
		return neutral_zone.QualityLow
	}

	if slices.Contains(playerNames, zone.Name) {
		return neutral_zone.QualityUnknown
	}

	switch zone_helpers.GetZoneTypeFromName(zone.Name) {
	case preview.ZoneTypeHub:
		return neutral_zone.QualityHighest
	case preview.ZoneTypePlayer:
		return neutral_zone.QualityUnknown
	case preview.ZoneTypeNeutral:
		return this.GetQuality(zone)
	case preview.ZoneTypeUnknown:
		fallthrough
	default:
		return neutral_zone.QualityUnknown
	}
}

func (this *ZoneClassifier) GetConnectionGuardQuality(
	zoneA, zoneB string,
	zones []entities.Zone,
	playerNames []string,
) neutral_zone.Quality {
	if zoneA != "" && zoneB != "" &&
		slices.Contains(playerNames, zoneA) && slices.Contains(playerNames, zoneB) {
		return neutral_zone.QualityUnknown
	}

	return max(
		this.GetGuardQuality(zoneA, zones, playerNames),
		this.GetGuardQuality(zoneB, zones, playerNames))
}

func (this *ZoneClassifier) getCenterQuality(zone entities.Zone) neutral_zone.Quality {
	if len(zone.GuardedContentPool) == 0 && len(zone.UnguardedContentPool) == 0 {
		return neutral_zone.QualityUnknown
	}

	isHighestTier := func(value string) bool { return strings.Contains(value, "_t5_") }
	if linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(value string) bool { return value == registry.GetResourcesContentPoolValues().TreasureZoneRich }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(isHighestTier) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(isHighestTier) {
		return neutral_zone.QualityHighest
	}

	return neutral_zone.QualityUnknown
}

func (this *ZoneClassifier) getTreasureQuality(zone entities.Zone) neutral_zone.Quality {
	if len(zone.GuardedContentPool) == 0 && len(zone.UnguardedContentPool) == 0 {
		return neutral_zone.QualityUnknown
	}

	isHighTier := func(value string) bool { return strings.Contains(value, "_t4_") || strings.Contains(value, "_t5_") }
	if linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(value string) bool { return value == registry.GetResourcesContentPoolValues().StartZoneRich }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(isHighTier) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(isHighTier) {
		return neutral_zone.QualityHigh
	}

	isMediumTier := func(value string) bool { return strings.Contains(value, "_t3_") }
	if linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(value string) bool { return value == registry.GetResourcesContentPoolValues().StartZoneMedium }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(isMediumTier) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(isMediumTier) {
		return neutral_zone.QualityMedium
	}

	return neutral_zone.QualityUnknown
}

func (this *ZoneClassifier) getSidesQuality(zone entities.Zone) neutral_zone.Quality {
	if len(zone.GuardedContentPool) == 0 && len(zone.UnguardedContentPool) == 0 {
		return neutral_zone.QualityUnknown
	}

	isLowTier := func(value string) bool { return strings.Contains(value, "_t2_") }
	if linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(value string) bool { return value == registry.GetResourcesContentPoolValues().StartZonePoor }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(isLowTier) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(isLowTier) {
		return neutral_zone.QualityLow
	}

	isLowestTier := func(value string) bool { return strings.Contains(value, "_t1_") }
	if linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(value string) bool { return value == registry.GetResourcesContentPoolValues().StartZoneVeryPoor }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(isLowestTier) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(isLowestTier) {
		return neutral_zone.QualityLowest
	}

	return neutral_zone.QualityUnknown
}
