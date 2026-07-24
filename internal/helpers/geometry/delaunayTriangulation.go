package geometry

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

// CreateDelaunayTriangulation returns normalized, deterministically ordered
// Delaunay edges for the supplied positions using Bowyer-Watson insertion.
func CreateDelaunayTriangulation(positions []data.Vec2[float64]) []data.Vec2[int] {
	count := len(positions)
	if count <= 1 {
		return nil
	}
	if count == 2 {
		return []data.Vec2[int]{data.NewVec2(0, 1)}
	}
	minimumPosition, maximumPosition := GetPositionBounds(positions)
	delta := max(maximumPosition.X-minimumPosition.X, maximumPosition.Y-minimumPosition.Y) * 10

	superPoints := make([]data.Vec2[float64], count+3)
	copy(superPoints, positions)
	superPoints[count] = data.NewVec2(minimumPosition.X-delta, minimumPosition.Y-delta*3)
	superPoints[count+1] = data.NewVec2(minimumPosition.X+delta*3, minimumPosition.Y-delta)
	superPoints[count+2] = data.NewVec2(minimumPosition.X, minimumPosition.Y+delta*3)

	triangleIndexes := []data.Vec3[int]{data.NewVec3(count, count+1, count+2)}
	for index := range positions {
		triangleIndexes = insertPointIntoTriangulation(triangleIndexes, superPoints, index)
	}

	return extractRealEdges(triangleIndexes, count)
}

func normalizeEdge(firstIndex, secondIndex int) data.Vec2[int] {
	if firstIndex > secondIndex {
		firstIndex, secondIndex = secondIndex, firstIndex
	}
	return data.NewVec2(firstIndex, secondIndex)
}

func insertPointIntoTriangulation(
	triangleIndexes []data.Vec3[int],
	superPoints []data.Vec2[float64],
	index int,
) []data.Vec3[int] {
	point := superPoints[index]
	edgeCount := map[data.Vec2[int]]int{}
	var kept []data.Vec3[int]
	for _, triangleIndex := range triangleIndexes {
		triangle := [3]data.Vec2[float64]{
			superPoints[triangleIndex.X],
			superPoints[triangleIndex.Y],
			superPoints[triangleIndex.Z],
		}
		if inCircumscribedCircle(triangle, point) {
			edgeCount[normalizeEdge(triangleIndex.X, triangleIndex.Y)]++
			edgeCount[normalizeEdge(triangleIndex.Y, triangleIndex.Z)]++
			edgeCount[normalizeEdge(triangleIndex.Z, triangleIndex.X)]++
		} else {
			kept = append(kept, triangleIndex)
		}
	}
	for edge, occurrences := range edgeCount {
		if occurrences == 1 {
			kept = append(kept, data.NewVec3(edge.X, edge.Y, index))
		}
	}
	return kept
}

func extractRealEdges(triangleIndexes []data.Vec3[int], count int) []data.Vec2[int] {
	edgeSet := map[data.Vec2[int]]bool{}
	for _, triangleIndex := range triangleIndexes {
		if triangleIndex.X < count && triangleIndex.Y < count && triangleIndex.Z < count {
			edgeSet[normalizeEdge(triangleIndex.X, triangleIndex.Y)] = true
			edgeSet[normalizeEdge(triangleIndex.Y, triangleIndex.Z)] = true
			edgeSet[normalizeEdge(triangleIndex.Z, triangleIndex.X)] = true
		}
	}
	result := make([]data.Vec2[int], 0, len(edgeSet))
	for edge := range edgeSet {
		result = append(result, edge)
	}
	slices.SortFunc(result, func(first, second data.Vec2[int]) int {
		if first.X != second.X {
			return first.X - second.X
		}
		return first.Y - second.Y
	})
	return result
}

// inCircumscribedCircle checks if "point" lies strictly inside the circumcircle of triangle ABC.
//
// The condition is satisfied if the following 4x4 determinant is positive:
//
//	| A_x      A_y      A_x²     + A_y²      1 |
//	| B_x      B_y      B_x²     + B_y²      1 |
//	| C_x      C_y      C_x²     + C_y²      1 | > 0
//	| point_x  point_y  point_x² + point_y²  1 |
//
// By subtracting the last row from the first three, it reduces to:
//
//	| A_x - point_x  A_y - point_y  (A_x - point_x)² + (A_y - point_y)² |
//	| B_x - point_x  B_y - point_y  (B_x - point_x)² + (B_y - point_y)² | > 0
//	| C_x - point_x  C_y - point_y  (C_x - point_x)² + (C_y - point_y)² |
//
// The circumscribed circle of a triangle is the unique circle that passes
// through all three vertices. Below: an equilateral triangle inscribed in its
// circumscribed circle (vertices on the circle, sides as chords).
//
//	```
//	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣠⠤⠤⠤⢤⡤⠤⠤⠤⣄⣀⡀
//	⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⠔⠚⠉⠀⠀⠀⠀⣰⠃⠘⣆⠀⠀⠀⠀⠉⠓⠢⣄⡀
//	⠀⠀⠀⠀⠀⠀⠀⣠⠖⠉⠀⠀⠀⠀⠀⠀⠀⡴⠁⠀⠀⠈⢦⠀⠀⠀⠀⠀⠀⠀⠉⠲⣄
//	⠀⠀⠀⠀⠀⣠⠞⠁⠀⠀⠀⠀⠀⠀⠀⠀⡼⠁⠀⠀⠀⠀⠈⢧⠀⠀⠀⠀⠀⠀⠀⠀⠈⠳⣄
//	⠀⠀⠀⢀⡜⠁⠀⠀⠀⠀⠀⠀⠀⠀⢀⡞⠀⠀⠀⠀⠀⠀⠀⠀⢳⡀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢣⡀
//	⠀⠀⢀⠞⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠳⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠳⡀
//	⠀⢀⡞⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⢳⡀
//	⠀⣸⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⣇
//	⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⡰⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⢆⠀⠀⠀⠀⠀⠀⠀⠀⢸
//	⠀⡇⠀⠀⠀⠀⠀⠀⠀⡼⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢧⠀⠀⠀⠀⠀⠀⠀⢸
//	⠀⡇⠀⠀⠀⠀⠀⢀⡜⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢣⡀⠀⠀⠀⠀⠀⢸
//	⠀⡇⠀⠀⠀⠀⢀⡞⠀⠀⠀⠀⠀. point⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢳⡀⠀⠀⠀⠀⢸
//	⠀⢹⠀⠀⠀⢠⠎⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠱⡄⠀⠀⠀⡏
//	⠀⠈⢧⠀⢠⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⡄⠀⡼
//	⠀⠀⠈⢶⣃⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣘⡶⠁
//	⠀⠀⠀⠈⢣⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡜⠁
//	⠀⠀⠀⠀⠀⠙⢦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡴⠋
//	⠀⠀⠀⠀⠀⠀⠀⠙⠦⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⠴⠋
//	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠢⢤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⡤⠔⠋⠁
//	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠙⠒⠒⠒⠒⠒⠒⠒⠒⠋⠉
//	```
func inCircumscribedCircle(triangle [3]data.Vec2[float64], point data.Vec2[float64]) bool {
	deltaA := triangle[0].Subtract(point)
	deltaB := triangle[1].Subtract(point)
	deltaC := triangle[2].Subtract(point)

	determinant := deltaA.X*(deltaB.Y*deltaC.SquaredLength()-deltaC.Y*deltaB.SquaredLength()) -
		deltaA.Y*(deltaB.X*deltaC.SquaredLength()-deltaC.X*deltaB.SquaredLength()) +
		deltaA.SquaredLength()*deltaB.CrossProduct(deltaC)

	// The determinant's sign flips with the triangle's winding order;
	// multiply by the orientation so the test is winding-independent.
	orientation := (deltaB.X-deltaA.X)*(deltaC.Y-deltaA.Y) - (deltaB.Y-deltaA.Y)*(deltaC.X-deltaA.X)
	return determinant*orientation > 0
}
