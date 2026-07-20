package neutral_zone

import (
	"cmp"
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
)

type Plan struct {
	Label       string
	Quality     Quality
	CastleCount int
}

func (this Plan) GetBalanceScore() float64 {
	return this.Quality.GetBalanceScore() + math.Min(float64(this.CastleCount), 4)*0.15
}

type Plans []Plan

func NewNeutralZonePlansSorted(plans Plans) *Plans {
	sortedPlans := make(Plans, len(plans))
	copy(sortedPlans, plans)
	sortedPlans.sort()
	return &sortedPlans
}

func NewNeutralZonePlansSortedByBalance(plans Plans) *Plans {
	sortedPlans := make(Plans, len(plans))
	copy(sortedPlans, plans)
	sortedPlans.sortByBalanceScoreDescending()
	return &sortedPlans
}

func (this *Plans) AddPlans(plans ...Plan) {
	*this = append(*this, plans...)
}

func (this *Plans) AddPlan(label string, quality Quality, castleCount int) {
	*this = append(*this, Plan{
		Label:       label,
		Quality:     quality,
		CastleCount: castleCount,
	})
}

func (this *Plans) GetQuality(label string) Quality {
	plan, ok := linq.FromSlice(*this).First(func(x Plan) bool { return x.Label == label })
	if !ok {
		return QualityUnknown
	}

	return plan.Quality
}

func (this *Plans) GetTier(label string) int {
	plan, ok := linq.FromSlice(*this).First(func(x Plan) bool { return x.Label == label })
	if !ok {
		return QualityUnknown.GetIndex()
	}

	return plan.Quality.GetIndex()
}

func (this *Plans) Any() bool {
	return len(*this) > 0
}

func (this *Plans) Swap(firstIndex, secondIndex int) {
	(*this)[firstIndex], (*this)[secondIndex] = (*this)[secondIndex], (*this)[firstIndex]
}

func (this *Plans) SortByBalanceScoreAscending() {
	slices.SortStableFunc(*this, func(a, b Plan) int {
		if comparison := cmp.Compare(a.GetBalanceScore(), b.GetBalanceScore()); comparison != 0 {
			return comparison
		}

		return cmp.Compare(a.Label, b.Label)
	})
}

func (this *Plans) sort() {
	slices.SortStableFunc(*this, func(a, b Plan) int {
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

func (this *Plans) sortByBalanceScoreDescending() {
	slices.SortStableFunc(*this, func(a, b Plan) int {
		if comparison := cmp.Compare(b.GetBalanceScore(), a.GetBalanceScore()); comparison != 0 {
			return comparison
		}

		return cmp.Compare(a.Label, b.Label)
	})
}
