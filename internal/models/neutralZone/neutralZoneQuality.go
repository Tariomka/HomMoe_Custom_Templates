package neutralZone

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// Quality is the tier of a neutral zone.
type Quality int

const (
	QualityLow    Quality = 0
	QualityMedium Quality = 1
	QualityHigh   Quality = 2
)

func (this Quality) GetGuardValue() int {
	switch this {
	case QualityHigh:
		return 25_000
	case QualityMedium:
		return 20_000
	case QualityLow:
		fallthrough
	default:
		return 15_000
	}
}

func GetQualityFrom(zone entities.Zone) Quality {
	pool := ""
	if len(zone.GuardedContentPool) > 0 {
		pool = zone.GuardedContentPool[0]
	}

	switch {
	case strings.Contains(pool, "_t4_") || strings.Contains(pool, "_t5_"):
		return QualityHigh
	case strings.Contains(pool, "_t1_") || strings.Contains(pool, "_t2_"):
		return QualityLow
	case strings.Contains(pool, "_t3_"):
		fallthrough
	default:
		return QualityMedium
	}
}
