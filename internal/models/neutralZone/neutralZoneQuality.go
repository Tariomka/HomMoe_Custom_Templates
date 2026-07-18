package neutralZone

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type Quality int8 // Tier of a neutral zone

const (
	QualityUnknown Quality = iota - 1 // Zone quality could not be determined, or is not a neutral zone
	QualityLowest                     // Plastic zone quality
	QualityLow                        // Bronze zone quality
	QualityMedium                     // Silver zone quality
	QualityHigh                       // Gold zone quality
	QualityHighest                    // Platinum zone quality
)

func (this Quality) GetGuardValue() int {
	switch this {
	case QualityHighest:
		return 30_000
	case QualityHigh:
		return 25_000
	case QualityMedium:
		return 20_000
	case QualityLow:
		return 15_000
	case QualityLowest:
		return 10_000
	case QualityUnknown:
		fallthrough
	default:
		return 0
	}
}

func (this Quality) GetBalanceScore() float64 {
	switch this {
	case QualityHighest:
		return 4.0
	case QualityHigh:
		return 3.0
	case QualityMedium:
		return 2.0
	case QualityLow:
		return 1.0
	case QualityLowest:
		return 0.5
	case QualityUnknown:
		fallthrough
	default:
		return 0
	}
}

func (this Quality) GetIndex() int {
	return int(this)
}

func GetQualityFrom(zone entities.Zone) Quality {
	// v2 of GetQualityFrom, assumes the Quality by the same values as Profile information
	// Will need to make it more robust in the future, if rmg.json files will be loaded

	if strings.HasPrefix(zone.Name, "Spawn-") {
		return QualityUnknown // Player spawn zones are not neutral zones
	}

	switch layoutValues := registry.GetLayoutValues(); zone.Layout {
	case layoutValues.Center:
		if quality := checkQualityForCenter(zone); quality != QualityUnknown {
			return quality
		}

	case layoutValues.TreasureZone:
		if quality := checkQualityForTreasure(zone); quality != QualityUnknown {
			return quality
		}

	case layoutValues.Sides:
		if quality := checkQualityForSides(zone); quality != QualityUnknown {
			return quality
		}
	}

	return QualityUnknown
}

func checkQualityForCenter(zone entities.Zone) Quality {
	if len(zone.GuardedContentPool) == 0 && len(zone.UnguardedContentPool) == 0 {
		return QualityUnknown
	}

	if highestCheck := func(x string) bool {
		return strings.Contains(x, "_t5_")
	}; linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(x string) bool { return x == registry.GetResourcesContentPoolValues().TreasureZoneRich }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(highestCheck) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(highestCheck) {
		return QualityHighest
	}

	return QualityUnknown
}

func checkQualityForTreasure(zone entities.Zone) Quality {
	if len(zone.GuardedContentPool) == 0 && len(zone.UnguardedContentPool) == 0 {
		return QualityUnknown
	}

	if highCheck := func(x string) bool {
		return strings.Contains(x, "_t4_") || strings.Contains(x, "_t5_")
	}; linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(x string) bool { return x == registry.GetResourcesContentPoolValues().StartZoneRich }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(highCheck) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(highCheck) {
		return QualityHigh
	}

	if mediumCheck := func(x string) bool {
		return strings.Contains(x, "_t3_")
	}; linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(x string) bool { return x == registry.GetResourcesContentPoolValues().StartZoneMedium }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(mediumCheck) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(mediumCheck) {
		return QualityMedium
	}

	return QualityUnknown
}

func checkQualityForSides(zone entities.Zone) Quality {
	if len(zone.GuardedContentPool) == 0 && len(zone.UnguardedContentPool) == 0 {
		return QualityUnknown
	}

	if lowCheck := func(x string) bool {
		return strings.Contains(x, "_t2_")
	}; linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(x string) bool { return x == registry.GetResourcesContentPoolValues().StartZonePoor }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(lowCheck) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(lowCheck) {
		return QualityLow
	}

	if lowestCheck := func(x string) bool {
		return strings.Contains(x, "_t1_")
	}; linq.FromSlice(zone.ResourcesContentPool).
		AllFunc(func(x string) bool { return x == registry.GetResourcesContentPoolValues().StartZoneVeryPoor }) ||
		linq.FromSlice(zone.GuardedContentPool).AllFunc(lowestCheck) ||
		linq.FromSlice(zone.UnguardedContentPool).AllFunc(lowestCheck) {
		return QualityLowest
	}

	return QualityUnknown
}
