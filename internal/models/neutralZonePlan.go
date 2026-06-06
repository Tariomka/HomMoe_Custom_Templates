package models

import (
	"cmp"
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
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

func NewNeutralZonePlansSorted(plans NeutralZonePlans) *NeutralZonePlans {
	sortedPlans := make(NeutralZonePlans, len(plans))
	copy(sortedPlans, plans)
	sortedPlans.sort()
	return &sortedPlans
}

func NewNeutralZonePlansSortedByBalance(plans NeutralZonePlans) *NeutralZonePlans {
	sortedPlans := make(NeutralZonePlans, len(plans))
	copy(sortedPlans, plans)
	sortedPlans.sortByBalanceScoreDescending()
	return &sortedPlans
}

func (this *NeutralZonePlans) AddPlans(plans ...NeutralZonePlan) {
	*this = append(*this, plans...)
}

func (this *NeutralZonePlans) AddPlan(label string, quality NeutralZoneQuality, castleCount int) {
	*this = append(*this, NeutralZonePlan{
		Label:       label,
		Quality:     quality,
		CastleCount: castleCount,
	})
}

func (this *NeutralZonePlans) AddMediumPlan(label string, castleCount int) {
	*this = append(*this, NeutralZonePlan{
		Label:       label,
		Quality:     QualityMedium,
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

func (this *NeutralZonePlans) GetTier(label string) int {
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

func (this *NeutralZonePlans) Any() bool {
	return len(*this) > 0
}

func (this *NeutralZonePlans) Swap(firstIndex, secondIndex int) {
	(*this)[firstIndex], (*this)[secondIndex] = (*this)[secondIndex], (*this)[firstIndex]
}

func (this *NeutralZonePlans) SortByBalanceScoreAscending() {
	slices.SortStableFunc(*this, func(a, b NeutralZonePlan) int {
		if comparison := cmp.Compare(a.GetBalanceScore(), b.GetBalanceScore()); comparison != 0 {
			return comparison
		}

		return cmp.Compare(a.Label, b.Label)
	})
}

func (this *NeutralZonePlans) sort() {
	slices.SortStableFunc(*this, func(a, b NeutralZonePlan) int {
		// a.Quality > b.Quality
		if comparison := cmp.Compare(b.Quality, a.Quality); comparison != 0 {
			return comparison
		}

		// a.CastleCount > b.CastleCount
		if comparison := cmp.Compare(b.CastleCount, a.CastleCount); comparison != 0 {
			return comparison
		}

		// a.Label < b.Label
		return cmp.Compare(a.Label, b.Label)
	})
}

func (this *NeutralZonePlans) sortByBalanceScoreDescending() {
	slices.SortStableFunc(*this, func(a, b NeutralZonePlan) int {
		if comparison := cmp.Compare(b.GetBalanceScore(), a.GetBalanceScore()); comparison != 0 {
			return comparison
		}

		return cmp.Compare(a.Label, b.Label)
	})
}
