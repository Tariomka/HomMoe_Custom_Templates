package models

import (
	"math"
	"slices"
)

type Vector2 struct {
	X, Y float64
}

type Positions []Vector2

func newPositions(size int) Positions {
	return make(Positions, size)
}

func CreatePositionsFromPlans(orderedLabels, playerLabels []string, neutralZonePlans NeutralZonePlans) Positions {
	count := len(orderedLabels)
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
	for i, l := range orderedLabels {
		tier := 0
		if !slices.Contains(playerLabels, l) {
			tier = neutralZonePlans.GetTier(l)
		}
		byTier[tier] = append(byTier[tier], i)
	}

	positions := newPositions(count)
	for tier, indices := range byTier {
		radius := tierRadius(tier)
		nn := len(indices)
		offset := float64(tier) * math.Pi / math.Max(1, float64(nn))
		for j, idx := range indices {
			angle := 2*math.Pi*float64(j)/float64(nn) + offset
			jitter := float64(j%3-1) * 0.008
			positions[idx] = Vector2{
				X: math.Max(0.05, math.Min(0.95, 0.5+math.Cos(angle+jitter)*radius)),
				Y: math.Max(0.05, math.Min(0.95, 0.5+math.Sin(angle+jitter)*radius)),
			}
		}
	}
	return positions
}

func (this *Positions) GetShortestDistanceIndex(adjacencyIndexes [][]int) (indexA, indexB int, ok bool) {
	indexA, indexB = -1, -1
	ok = false
	if len(adjacencyIndexes) <= 1 {
		return
	}

	mainComp := map[int]bool{}
	for _, idx := range adjacencyIndexes[0] {
		mainComp[idx] = true
	}
	bestDist := math.MaxFloat64
	for _, a := range adjacencyIndexes[0] {
		for ci := 1; ci < len(adjacencyIndexes); ci++ {
			for _, b := range adjacencyIndexes[ci] {
				dx := (*this)[a].X - (*this)[b].X
				dy := (*this)[a].Y - (*this)[b].Y
				d := dx*dx + dy*dy
				if d < bestDist {
					bestDist = d
					indexA, indexB = a, b
					ok = true
				}
			}
		}
	}
	return indexA, indexB, ok
}
