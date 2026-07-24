package geometry

import (
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

func FindClosestAcrossComponents(
	positions []data.Vec2[float64],
	componentIndexes [][]int,
) (indexes data.Vec2[int], found bool) {
	indexes = data.NewVec2(-1, -1)
	if len(componentIndexes) <= 1 {
		return indexes, false
	}

	bestDistance := math.MaxFloat64
	for _, startIndex := range componentIndexes[0] {
		for _, indexesInComponent := range componentIndexes[1:] {
			for _, endIndex := range indexesInComponent {
				delta := positions[startIndex].Subtract(positions[endIndex])
				distance := delta.SquaredLength()
				if distance < bestDistance {
					bestDistance = distance
					indexes = data.NewVec2(startIndex, endIndex)
					found = true
				}
			}
		}
	}
	return indexes, found
}
