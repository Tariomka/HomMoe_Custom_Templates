package models

import (
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/linq"
)

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

type NeutralZonePlans []NeutralZonePlan

func (this *NeutralZonePlans) AddPlan(label string, quality NeutralZoneQuality, castleCount int) {
	*this = append(*this, NeutralZonePlan{
		Label:       label,
		Quality:     quality,
		CastleCount: castleCount,
	})
}

func (this *NeutralZonePlans) GetQuality(label string) NeutralZoneQuality {
	if len(*this) == 0 {
		return QualityMedium
	}
	plan, ok := linq.FromSlice(*this).First(func(x NeutralZonePlan) bool { return x.Label == label })
	if !ok {
		return QualityMedium
	}
	return plan.Quality
}

func (this *NeutralZonePlans) GetTier(label string, playerLabels []string) int {
	if slices.Contains(playerLabels, label) {
		return 0
	}
	plan, ok := linq.FromSlice(*this).First(func(x NeutralZonePlan) bool { return x.Label == label })
	if !ok {
		return 1
	}
	switch plan.Quality {
	case QualityHigh:
		return 3
	case QualityMedium:
		return 2
	default:
		return 1
	}
}
