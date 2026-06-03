package models

import "math"

type NeutralZonePlan struct {
	Label       string
	Quality     NeutralZoneQuality
	CastleCount int
}

func (this NeutralZonePlan) GetBalanceScore() float64 {
	var qualityScore float64
	switch this.Quality {
	case QualityHigh:
		qualityScore = 3.0
	case QualityMedium:
		qualityScore = 2.0
	default:
		qualityScore = 1.0
	}

	return qualityScore + math.Min(float64(this.CastleCount), 4)*0.15
}
