package models

import "math"

type Measurement [2]float64

type Measurements []Measurement

func newMeasurements(size int) Measurements {
	return make(Measurements, size)
}

func CreateMeasurementsFromPlans(orderedLetters, playerLetters []string, neutralZonePlans NeutralZonePlans) Measurements {
	count := len(orderedLetters)
	if count == 0 {
		return nil
	}

	tierRadius := func(tier int) float64 {
		switch tier {
		case 0:
			return 0.38
		case 1:
			return 0.27
		case 2:
			return 0.16
		default:
			return 0.06
		}
	}

	byTier := map[int][]int{}
	for i, l := range orderedLetters {
		tier := neutralZonePlans.GetTier(l, playerLetters)
		byTier[tier] = append(byTier[tier], i)
	}

	positions := newMeasurements(count)
	for tier, indices := range byTier {
		radius := tierRadius(tier)
		nn := len(indices)
		offset := float64(tier) * math.Pi / math.Max(1, float64(nn))
		for j, idx := range indices {
			angle := 2*math.Pi*float64(j)/float64(nn) + offset
			jitter := float64(j%3-1) * 0.008
			positions[idx] = [2]float64{
				math.Max(0.05, math.Min(0.95, 0.5+math.Cos(angle+jitter)*radius)),
				math.Max(0.05, math.Min(0.95, 0.5+math.Sin(angle+jitter)*radius)),
			}
		}
	}
	return positions
}
