package utils

import (
	"math"

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
	for i := range remaining {
		gap := int(math.Floor((float64(i) + 0.5) * float64(gapCount) / float64(remaining)))
		capacities[helpers.Clamp(gap, 0, gapCount-1)]++
	}
	return capacities
}

func AssignNeutralZonesToGaps(
	neutralZones models.NeutralZonePlans,
	capacities []int,
	preferInterior bool) []models.NeutralZonePlans {
	gaps := make([]models.NeutralZonePlans, len(capacities))
	loads := make([]float64, len(capacities))
	sortedZones := models.NewNeutralZonePlansSortedByBalance(neutralZones)
	for _, zonePlan := range *sortedZones {
		var candidates []int
		for i := range capacities {
			if len(gaps[i]) < capacities[i] {
				candidates = append(candidates, i)
			}
		}
		if len(candidates) == 0 {
			break
		}
		if preferInterior {
			var interior []int
			for _, c := range candidates {
				if c > 0 && c < len(capacities)-1 {
					interior = append(interior, c)
				}
			}
			if len(interior) > 0 {
				candidates = interior
			}
		}
		best := candidates[0]
		for _, candidate := range candidates[1:] {
			if loads[candidate] < loads[best] ||
				(loads[candidate] == loads[best] && len(gaps[candidate]) < len(gaps[best])) ||
				(loads[candidate] == loads[best] && len(gaps[candidate]) == len(gaps[best]) && candidate < best) {
				best = candidate
			}
		}
		gaps[best] = append(gaps[best], zonePlan)
		loads[best] += zonePlan.GetBalanceScore()
	}
	return gaps
}

func OrderNeutralsWithinGap(neutralZones models.NeutralZonePlans) models.NeutralZonePlans {
	if len(neutralZones) <= 1 {
		zones := models.NeutralZonePlans{}
		zones.AddPlans(neutralZones...)
		return zones
	}

	sortedZones := models.NewNeutralZonePlansSortedByBalance(neutralZones)
	slots := make(models.NeutralZonePlans, len(*sortedZones))
	lowIndex, highIndex := 0, len(*sortedZones)-1
	for i, zonePlan := range *sortedZones {
		if i%2 == 0 {
			slots[lowIndex] = zonePlan
			lowIndex++
		} else {
			slots[highIndex] = zonePlan
			highIndex--
		}
	}
	return slots
}

func OrderEdgeGap(neutralZones models.NeutralZonePlans, playerAtEnd bool) models.NeutralZonePlans {
	sorted := models.NewNeutralZonePlansSortedByBalance(neutralZones)
	if playerAtEnd {
		for i, j := 0, len(*sorted)-1; i < j; i, j = i+1, j-1 {
			sorted.Swap(i, j)
		}
	}
	return *sorted
}
