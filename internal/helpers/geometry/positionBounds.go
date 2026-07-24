package geometry

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"

func GetPositionBounds(positions []data.Vec2[float64]) (minimumPosition, maximumPosition data.Vec2[float64]) {
	minimumPosition = positions[0]
	maximumPosition = positions[0]
	for _, position := range positions[1:] {
		minimumPosition.X = min(minimumPosition.X, position.X)
		minimumPosition.Y = min(minimumPosition.Y, position.Y)
		maximumPosition.X = max(maximumPosition.X, position.X)
		maximumPosition.Y = max(maximumPosition.Y, position.Y)
	}
	return minimumPosition, maximumPosition
}
