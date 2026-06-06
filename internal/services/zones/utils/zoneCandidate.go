package utils

import (
	"cmp"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type candidate struct {
	letter    string
	minDist   int
	variance  float64
	quality   int
	hasCastle int
}

type hubZoneCandidates []candidate

func CreateHubZoneCandidates(neutralZones models.NeutralZonePlans, distancesByPlayer []map[string]int) *hubZoneCandidates {
	var candidates hubZoneCandidates
	for _, plan := range neutralZones {
		var dists []int
		for _, d := range distancesByPlayer {
			v, ok := d[plan.Label]
			if !ok {
				v = 999999
			}
			dists = append(dists, v)
		}
		minD := dists[0]
		sum := 0
		for _, d := range dists {
			if d < minD {
				minD = d
			}
			sum += d
		}
		mean := float64(sum) / float64(len(dists))
		variance := 0.0
		for _, d := range dists {
			diff := float64(d) - mean
			variance += diff * diff
		}
		variance /= float64(len(dists))
		hc := 0
		if plan.CastleCount > 0 {
			hc = 1
		}
		candidates = append(candidates, candidate{plan.Label, minD, variance, int(plan.Quality), hc})
	}
	return &candidates
}

func (this *hubZoneCandidates) SortForHubCity() *hubZoneCandidates {
	slices.SortStableFunc(*this, func(a, b candidate) int {
		// a.minDist > b.minDist
		if comparison := cmp.Compare(b.minDist, a.minDist); comparison != 0 {
			return comparison
		}

		// a.variance < b.variance
		if comparison := cmp.Compare(a.variance, b.variance); comparison != 0 {
			return comparison
		}

		// a.quality > b.quality
		if comparison := cmp.Compare(b.quality, a.quality); comparison != 0 {
			return comparison
		}

		// a.hasCastle > b.hasCastle
		return cmp.Compare(b.hasCastle, a.hasCastle)
	})

	return this
}

func (this *hubZoneCandidates) GetFirstCandidateLabel() string {
	if len(*this) == 0 {
		return ""
	}
	return (*this)[0].letter
}
