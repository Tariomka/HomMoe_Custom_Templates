package models

// NeutralZoneQuality is the tier of a neutral zone.
type NeutralZoneQuality int

const (
	QualityLow    NeutralZoneQuality = 0
	QualityMedium NeutralZoneQuality = 1
	QualityHigh   NeutralZoneQuality = 2
)

func (this NeutralZoneQuality) GetGuardValue() int {
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
