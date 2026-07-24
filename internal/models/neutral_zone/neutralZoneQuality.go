package neutral_zone

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
		return 35_000
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
		return 30_000
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

func (this Quality) GetName() string {
	switch this {
	case QualityHighest:
		return "Platinum"
	case QualityHigh:
		return "Gold"
	case QualityMedium:
		return "Silver"
	case QualityLow:
		return "Bronze"
	case QualityLowest:
		return "Plastic"
	case QualityUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

func GetQualityFromIndex(index int) Quality {
	switch index {
	case int(QualityLowest):
		return QualityLowest
	case int(QualityLow):
		return QualityLow
	case int(QualityMedium):
		return QualityMedium
	case int(QualityHigh):
		return QualityHigh
	case int(QualityHighest):
		return QualityHighest
	default:
		return QualityUnknown
	}
}

func GetQualityNames() []string {
	return []string{
		QualityLowest.GetName(),
		QualityLow.GetName(),
		QualityMedium.GetName(),
		QualityHigh.GetName(),
		// QualityHighest.GetLabel(), // Platinum quality is not allowed to be used as neutral zone quality for now
	}
}
