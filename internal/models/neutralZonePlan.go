package models

import (
	"math"
	"sort"

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
	sortedPlans.sortByBalanceScore()
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

func (this *NeutralZonePlans) sort() {
	sort.SliceStable(*this, func(i, j int) bool {
		if (*this)[i].Quality != (*this)[j].Quality {
			return (*this)[i].Quality > (*this)[j].Quality
		}
		if (*this)[i].CastleCount != (*this)[j].CastleCount {
			return (*this)[i].CastleCount > (*this)[j].CastleCount
		}
		return (*this)[i].Label < (*this)[j].Label
	})
}

func (this *NeutralZonePlans) sortByBalanceScore() {
	sort.SliceStable(*this, func(i, j int) bool {
		scoreI, scoreJ := (*this)[i].GetBalanceScore(), (*this)[j].GetBalanceScore()
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return (*this)[i].Label < (*this)[j].Label
	})
}
