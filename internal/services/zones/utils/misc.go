package utils

import (
	"math"
	"sort"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func GetEvenGapCapacities(gapCount, itemCount, minimumPerGap int) []int {
	if gapCount <= 0 {
		return nil
	}
	capacities := make([]int, gapCount)
	if itemCount <= 0 {
		return capacities
	}
	minimum := max(0, minimumPerGap)
	reserved := minimum * gapCount
	remaining := itemCount
	if minimum > 0 && itemCount >= reserved {
		for i := range capacities {
			capacities[i] = minimum
		}
		remaining -= reserved
	}
	for i := 0; i < remaining; i++ {
		gap := int(math.Floor((float64(i) + 0.5) * float64(gapCount) / float64(remaining)))
		capacities[helpers.Clamp(gap, 0, gapCount-1)]++
	}
	return capacities
}

func AssignNeutralZonesToGaps(neutralZones []models.NeutralZonePlan, caps []int, preferInterior bool) [][]models.NeutralZonePlan {
	gaps := make([][]models.NeutralZonePlan, len(caps))
	loads := make([]float64, len(caps))
	sorted := sortNeutralZonePlans(neutralZones)
	for _, nz := range sorted {
		var candidates []int
		for i := range caps {
			if len(gaps[i]) < caps[i] {
				candidates = append(candidates, i)
			}
		}
		if len(candidates) == 0 {
			break
		}
		if preferInterior {
			var interior []int
			for _, c := range candidates {
				if c > 0 && c < len(caps)-1 {
					interior = append(interior, c)
				}
			}
			if len(interior) > 0 {
				candidates = interior
			}
		}
		best := candidates[0]
		for _, c := range candidates[1:] {
			if loads[c] < loads[best] || (loads[c] == loads[best] && len(gaps[c]) < len(gaps[best])) || (loads[c] == loads[best] && len(gaps[c]) == len(gaps[best]) && c < best) {
				best = c
			}
		}
		gaps[best] = append(gaps[best], nz)
		loads[best] += CalculateBalanceScore(nz)
	}
	return gaps
}

func OrderNeutralsWithinGap(neutralZones []models.NeutralZonePlan) []models.NeutralZonePlan {
	if len(neutralZones) <= 1 {
		return append([]models.NeutralZonePlan{}, neutralZones...)
	}

	sorted := sortNeutralZonePlans(neutralZones)
	slots := make([]models.NeutralZonePlan, len(sorted))
	lo, hi := 0, len(sorted)-1
	for i, nz := range sorted {
		if i%2 == 0 {
			slots[lo] = nz
			lo++
		} else {
			slots[hi] = nz
			hi--
		}
	}
	return slots
}

func OrderEdgeGap(neutralZones []models.NeutralZonePlan, playerAtEnd bool) []models.NeutralZonePlan {
	sorted := sortNeutralZonePlans(neutralZones)
	if playerAtEnd {
		for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
			sorted[i], sorted[j] = sorted[j], sorted[i]
		}
	}
	return sorted
}

func CalculateBalanceScore(zone models.NeutralZonePlan) float64 {
	q := 1.0
	switch zone.Quality {
	case models.QualityHigh:
		q = 3.0
	case models.QualityMedium:
		q = 2.0
	}
	return q + math.Min(float64(zone.CastleCount), 4)*0.15
}

func sortNeutralZonePlans(neutralZones []models.NeutralZonePlan) []models.NeutralZonePlan {
	sorted := make([]models.NeutralZonePlan, len(neutralZones))
	copy(sorted, neutralZones)
	sort.SliceStable(sorted, func(i, j int) bool {
		scoreI, scoreJ := CalculateBalanceScore(sorted[i]), CalculateBalanceScore(sorted[j])
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return sorted[i].Label < sorted[j].Label
	})
	return sorted
}
