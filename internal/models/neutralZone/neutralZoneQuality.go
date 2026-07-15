package neutralZone

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// Quality is the tier of a neutral zone.
type Quality int

const (
	QualityLowest  Quality = 0
	QualityLow     Quality = 1
	QualityMedium  Quality = 2
	QualityHigh    Quality = 3
	QualityHighest Quality = 4
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
		fallthrough
	default:
		return 10_000
	}
}

func GetQualityFrom(zone entities.Zone) Quality {
	pool := ""
	if len(zone.GuardedContentPool) > 0 {
		pool = zone.GuardedContentPool[0]
	}
	resourcesPool := ""
	if len(zone.ResourcesContentPool) > 0 {
		resourcesPool = zone.ResourcesContentPool[0]
	}

	switch {
	case strings.Contains(pool, "_t5_") && strings.Contains(resourcesPool, "treasure_zone_rich"):
		return QualityHighest
	case strings.Contains(pool, "_t4_") || strings.Contains(pool, "_t5_"):
		return QualityHigh
	case strings.Contains(pool, "_t1_"):
		return QualityLowest
	case strings.Contains(pool, "_t2_"):
		return QualityLow
	case strings.Contains(pool, "_t3_"):
		fallthrough
	default:
		return QualityMedium
	}
}
