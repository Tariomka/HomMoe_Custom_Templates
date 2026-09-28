package connectionCurveLayout_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// fixtureZoneRadius makes the obstacle clearance 46px and a detour's reach 52px,
// so every expected coordinate below is an exact integer.
const fixtureZoneRadius = 38.0

// chordPositions places A and B on a horizontal chord with nothing between them.
// Its positive normal points up, towards smaller y.
func chordPositions() map[string]data.Vec2[float64] {
	return map[string]data.Vec2[float64]{
		"A": data.NewVec2(140.0, 350.0),
		"B": data.NewVec2(560.0, 350.0),
	}
}

// chordWithObstacle adds a zone D at the given point to the A-B chord.
func chordWithObstacle(x, y float64) map[string]data.Vec2[float64] {
	positions := chordPositions()
	positions["D"] = data.NewVec2(x, y)

	return positions
}

func newLayout(positions map[string]data.Vec2[float64]) preview.ConnectionCurveLayout {
	return preview.ConnectionCurveLayout{Positions: positions, ZoneRadius: fixtureZoneRadius}
}

func newConnection(from, to string) template_model.Connection {
	return template_model.Connection{From: from, To: to}
}

// parallelConnections returns count connections from A to B.
func parallelConnections(count int) []template_model.Connection {
	connections := make([]template_model.Connection, 0, count)
	for range count {
		connections = append(connections, newConnection("A", "B"))
	}

	return connections
}

func controlPoints(curves []preview.ConnectionCurve) []data.Vec2[float64] {
	points := make([]data.Vec2[float64], 0, len(curves))
	for _, curve := range curves {
		points = append(points, curve.Control)
	}

	return points
}

func connectionIndices(curves []preview.ConnectionCurve) []int {
	indices := make([]int, 0, len(curves))
	for _, curve := range curves {
		indices = append(indices, curve.ConnectionIndex)
	}

	return indices
}
