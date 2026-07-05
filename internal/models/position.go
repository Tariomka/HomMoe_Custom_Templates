package models

import (
	"math"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
)

type Position = data.Vec2[float64]
type ConnectionIndexes = data.Vec2[int]

type Positions []Position

func CreatePositionsFromPlans(orderedLabels, playerLabels []string, neutralZonePlans NeutralZonePlans) Positions {
	count := len(orderedLabels)
	if count == 0 {
		return nil
	}

	getTierRadius := func(tier int) float64 {
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

	byTier := make(map[int][]int)
	for i, label := range orderedLabels {
		tier := 0
		if !slices.Contains(playerLabels, label) {
			tier = neutralZonePlans.GetTier(label)
		}
		byTier[tier] = append(byTier[tier], i)
	}

	positions := make(Positions, count)
	for tier, indices := range byTier {
		radius := getTierRadius(tier)
		nn := float64(len(indices))
		offset := float64(tier) * math.Pi / nn
		for i, idx := range indices {
			angle := 2*math.Pi*float64(i)/nn + offset
			jitter := float64(i%3-1) * 0.008
			positions[idx] = data.NewVec2(
				max(0.05, min(0.95, 0.5+math.Cos(angle+jitter)*radius)),
				max(0.05, min(0.95, 0.5+math.Sin(angle+jitter)*radius)))
		}
	}
	return positions
}

func (this *Positions) Add(position Position) { *this = append(*this, position) }

func (this *Positions) GetShortestDistanceIndex(adjacencyIndexes [][]int) (indices ConnectionIndexes, ok bool) {
	indices = data.NewVec2(-1, -1)
	if len(adjacencyIndexes) <= 1 {
		return
	}

	bestDist := math.MaxFloat64
	for _, start := range adjacencyIndexes[0] {
		for _, indexes := range adjacencyIndexes[1:] {
			for _, end := range indexes {
				delta := (*this)[start].Subtract((*this)[end])
				distance := delta.SquaredLength()
				if distance < bestDist {
					bestDist = distance
					indices = data.NewVec2(start, end)
					ok = true
				}
			}
		}
	}
	return indices, ok
}

// CreateDelaunayTriangulation creates a Delaunay triangulation of the positions,
// which divides the positions into a mesh of contiguous, non-overlapping triangles,
// and returns the edges as pairs of indexes. Delaunay triangulation ensures that no point
// falls inside the circumscribed circle (the circle that touches all three vertices) of any triangle.
func (this *Positions) CreateDelaunayTriangulation() []ConnectionIndexes {
	count := len(*this)
	if count <= 1 {
		return nil
	}
	if count == 2 {
		return []ConnectionIndexes{data.NewVec2(0, 1)}
	}
	minPos, maxPos := this.GetMinAndMaxPositions()
	delta := max(maxPos.X-minPos.X, maxPos.Y-minPos.Y) * 10

	superPts := make([]Position, count+3)
	copy(superPts, *this)
	superPts[count] = data.NewVec2(minPos.X-delta, minPos.Y-delta*3)
	superPts[count+1] = data.NewVec2(minPos.X+delta*3, minPos.Y-delta)
	superPts[count+2] = data.NewVec2(minPos.X, minPos.Y+delta*3)

	trianglesIndexes := []data.Vec3[int]{data.NewVec3(count, count+1, count+2)}

	normalizeEdge := func(a, b int) ConnectionIndexes {
		if a > b {
			a, b = b, a
		}
		return data.NewVec2(a, b)
	}

	for index := range *this {
		point := superPts[index]
		// Split triangles into "bad" ones (circumscribed circle contains the
		// point) and kept ones. Edges of the bad region that appear exactly
		// once form its boundary; re-triangulate the cavity from those.
		edgeCount := map[ConnectionIndexes]int{}
		var kept []data.Vec3[int]
		for _, triangleIndexes := range trianglesIndexes {
			triangle := [3]Position{
				superPts[triangleIndexes.X],
				superPts[triangleIndexes.Y],
				superPts[triangleIndexes.Z],
			}
			if inCircumscribedCircle(triangle, point) {
				edgeCount[normalizeEdge(triangleIndexes.X, triangleIndexes.Y)]++
				edgeCount[normalizeEdge(triangleIndexes.Y, triangleIndexes.Z)]++
				edgeCount[normalizeEdge(triangleIndexes.Z, triangleIndexes.X)]++
			} else {
				kept = append(kept, triangleIndexes)
			}
		}
		for e, occurrences := range edgeCount {
			if occurrences == 1 {
				kept = append(kept, data.NewVec3(e.X, e.Y, index))
			}
		}
		trianglesIndexes = kept
	}

	edgeSet := map[ConnectionIndexes]bool{}
	for _, t := range trianglesIndexes {
		if t.X < count && t.Y < count && t.Z < count {
			edgeSet[normalizeEdge(t.X, t.Y)] = true
			edgeSet[normalizeEdge(t.Y, t.Z)] = true
			edgeSet[normalizeEdge(t.Z, t.X)] = true
		}
	}
	result := make([]ConnectionIndexes, 0, len(edgeSet))
	for e := range edgeSet {
		result = append(result, e)
	}
	return result
}

func (this *Positions) GetMinAndMaxPositions() (minPos, maxPos Position) {
	minPos = (*this)[0]
	maxPos = (*this)[0]
	for _, position := range (*this)[1:] {
		minPos.X = min(minPos.X, position.X)
		minPos.Y = min(minPos.Y, position.Y)
		maxPos.X = max(maxPos.X, position.X)
		maxPos.Y = max(maxPos.Y, position.Y)
	}
	return
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
func inCircumscribedCircle(triangle [3]Position, point Position) bool {
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

// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⣖⡒⠲⢤⡀⠀⡞⠹⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣦⠀⠉⠲⢿⣶⣬⡷⣶⠤⢤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡶⢼⣯⣳⣦⡀⠀⡝⠁⠀⢹⠀⠈⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢳⡀⠀⠀⡈⠁⠰⢇⣀⣀⡼⣄⣠⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⠟⣿⣋⣡⠄⣴⠞⢦⡀⠀⠈⠉⣸⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⠏⢰⠁⠈⢹⣿⠿⠷⢆⣈⣀⣀⣰⣏⣀⣀⠀⠀⠀⢰⢿⡇⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⠋⢠⠏⠑⢴⠋⠁⠀⠀⠀⠀⠀⠀⠈⠙⠲⣄⠉⠲⢄⣸⣄⠓⢤⣀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⡏⠀⡏⡀⠀⡏⠀⠀⠀⠀⠀⠀⢀⡤⠖⢶⡚⢻⣧⣀⡀⠉⢿⠿⠿⣾⠓⢦⡀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⠀⢸⠁⠙⠒⡇⠀⠀⠀⠀⢀⡞⠁⠀⢀⢸⡉⠚⠛⠛⠣⢰⠃⠀⠀⣰⠀⠀⡇⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⡶⢶⡦⢤⡀⢸⡇⢠⠦⡀⠀⠀⡇⠀⠀⢀⡴⢋⡤⠖⡋⠉⢉⣑⣦⡤⠖⠀⣈⠷⣄⠒⠙⠲⠞⢻⡄
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠈⠛⡿⠛⠃⢻⣼⠀⡼⠀⠉⠓⠲⡇⠀⣠⠏⢠⠉⠀⢀⡽⠛⠉⠉⠙⠓⠾⠭⠶⣆⣈⡁⠀⠀⣀⡾⠁
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⠀⡤⠴⠖⠀⢸⡏⢠⠟⢦⣀⠀⠀⣇⡴⠃⡴⠀⠒⣶⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⠀⣇⠀⣀⡀⢸⡇⠀⠀⠀⠈⠉⠁⡟⢀⡾⣀⣀⣰⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⠞⠋⠙⢦⠁⠀⣼⣧⠀⠀⠀⠀⠀⠀⣿⡼⠀⠀⢩⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡼⢿⡏⠀⠀⠀⠘⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠳⢤⡄⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣀⣤⠤⠖⠚⡉⢀⣾⠀⠀⠀⢀⣀⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠃⠰⣦⣀⠂⠌⡐⢀⢒⣧⠀⢀⣤⠞⣽⠁⠐⠶⠤⣄⠀⠀⠀⠀⢀⣀⣀⠀⣼⠀⣿⠉⠓⠲⠦⠤⠴⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⠃⡌⢐⠠⠙⢷⡄⡈⠄⣺⠿⣦⠞⠁⠀⡏⠀⠀⠀⠀⠀⠙⠦⠞⠉⠉⠉⠀⢰⣇⠀⣸⠇⠨⠐⣤⠞⠡⢸⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⠀⣸⡏⠰⢀⠆⡀⠇⡀⢹⣆⢸⡿⠎⠁⠀⠀⢰⡇⠀⠀⢰⡀⠀⠀⠀⠀⠀⠀⠀⠀⣾⠉⣷⣿⠁⠆⣿⠇⠈⡰⠀⢿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠀⢠⡟⠠⢁⠂⠰⠐⡀⠒⠠⢉⡾⠁⠀⠀⠀⠀⣼⠁⠀⠀⠀⢱⡀⠀⠀⠀⣠⠖⠀⣰⠇⠀⠈⠻⣇⢸⠇⡀⠃⠄⡁⢊⢷⣂⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⢠⠞⠠⢁⠂⠌⠠⠁⡄⢁⣲⡟⠁⠀⠀⠀⠀⣴⠃⠀⠀⠀⠀⠀⠙⠀⠒⠋⠁⠀⠀⣿⡀⠀⠀⠀⢻⡆⠐⡠⢁⠂⠔⡀⠪⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⣰⠋⠠⢁⠂⠌⠠⠁⠆⡐⣰⣏⡀⠀⠀⠀⠀⣰⠏⠀⠀⠀⠀⠹⣄⠀⠀⠀⣀⠀⠀⣸⢻⣧⠀⠀⠀⠀⢻⡄⠰⠀⢌⠐⠠⢁⠸⣷⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⢀⡼⠁⡌⢁⠂⢌⣠⡷⠞⠛⢛⣿⠧⠤⠄⠀⠀⣾⣁⣀⠀⠀⠀⠀⠀⠉⠀⠛⠋⠁⠀⣰⢯⠀⢹⣆⠀⠀⠀⠀⢿⣆⣉⠄⢊⢁⠂⠄⠹⣧⡀⠀⠀⠀⠀⠀⠀
// ⠀⠀⣠⠞⠠⢁⣰⡴⠞⢋⠁⠄⡈⠆⡀⢿⣆⠀⠀⠀⠀⢹⡄⠉⠛⠆⠀⠀⠀⠀⠀⢀⠀⠀⣰⡿⠃⠀⠀⠻⣦⠀⠀⢀⣼⡇⢉⠙⠳⢦⣌⡄⢁⠻⣕⡀⠀⠀⠀⠀⠀
// ⣤⡾⠃⣀⢲⠟⢉⠠⠈⠄⠨⠐⠄⠒⢠⠈⢽⣆⠀⠀⠀⢨⡇⠀⠀⠀⠸⡧⠄⡷⠖⠋⠀⣼⣿⡅⠀⠀⠀⠀⣿⠃⠀⠀⢈⡟⠠⠘⡀⠆⡈⠙⠷⣄⢻⣥⡀⠀⠀⠀⠀
// ⠙⢦⡐⠠⢂⠐⡀⢂⠡⢈⠡⠈⠌⠰⠀⢌⠠⢹⣆⠀⠀⢸⠀⠀⠀⠀⢸⡇⠀⠀⢀⣴⣽⠋⠀⢳⡀⠀⠀⠀⣿⠀⠀⠀⢸⡏⠠⠡⠐⡠⢈⠢⠐⡀⢂⣹⡧⠀⠀⠀⠀
// ⠀⠈⠻⡄⢂⠒⡀⢂⡐⠂⠄⠃⠌⢂⠉⠄⠂⣄⣹⣆⠀⢸⡆⠀⠀⢀⣾⠀⠀⣠⡿⠋⠀⠀⠀⠀⢻⡄⠀⠀⣿⠀⠀⠀⣿⠀⣁⠢⢁⠐⡠⢂⠡⠐⣰⡿⠁⠀⠀⠀⠀
// ⠀⠀⠀⠹⡆⠐⡐⠠⢀⠡⢊⠐⡈⠄⣈⣴⠷⠋⡁⢹⡆⠀⣧⠀⠀⢰⣇⣤⠾⠋⠀⠀⠀⠀⠀⠀⠀⢻⡄⠀⢹⡆⠀⢀⣿⡄⠀⠆⡐⢂⠐⡀⢂⢡⢿⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⢹⡄⠁⠆⡁⠢⠄⠌⣐⡴⢋⠁⡐⠐⣀⠂⣻⠀⠸⣧⣠⡟⠛⠛⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⡀⠀⣷⠀⣸⠃⢿⡈⢐⠀⠆⠨⠐⣠⣿⠃⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⢷⠉⡐⠠⠁⠌⣰⠟⡀⢂⠡⠀⡅⡀⠂⣿⠀⠀⠀⠙⢦⡀⠀⠘⣻⠄⠀⠀⠀⠠⣄⡀⠐⠚⣧⣴⡟⠀⣿⠀⡘⣧⠠⢈⢂⠁⢢⣿⠃⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⢸⡂⠄⡡⢈⣼⠃⡐⠠⢁⠂⠡⠐⢠⣡⣿⠀⠀⠘⣷⣌⢷⣀⡞⠁⠀⠀⠀⠀⠀⢀⣉⣭⠿⢯⠏⡀⠐⣿⣧⡀⢹⡆⠐⡠⢈⣼⠇⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠈⡇⠰⠀⣾⠇⠰⢀⠁⠆⡈⢱⣾⠿⠉⢹⡀⡆⡆⢿⣿⣾⠁⠿⣆⠀⠀⠀⢰⡿⠏⠁⠀⠀⣾⢸⣿⢰⣿⣿⣷⠀⣿⠀⠰⢀⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⣇⠂⣱⡏⠠⢁⠂⢌⣰⠾⢋⣿⠀⠀⠹⣧⣷⣹⡈⣿⡄⠀⠀⠈⠳⣄⠀⠐⣧⠀⠀⠀⠀⠙⠫⡿⢸⣿⡟⠻⣧⠙⠠⢁⢺⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⣧⠀⠻⢀⠁⣂⣾⡿⡿⠀⠸⣿⠀⠀⠀⠈⠹⣿⣷⣘⡧⠀⠀⠀⠀⠘⢷⡄⢻⡇⠀⠀⠀⠀⠀⣯⡼⠟⠀⠀⠹⣧⡁⠂⣼⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⣟⠠⢁⢂⣵⠟⢠⡞⠁⠀⠀⠉⠓⠒⣚⠉⠙⣆⠈⠀⠀⠀⠀⠀⠀⠀⠀⢹⠄⠙⠒⠶⢶⣚⠛⢲⡀⠀⠀⠀⠀⠘⣧⠀⣯⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⣯⠐⣨⠟⠁⠀⠀⠓⠲⢤⡀⠀⢦⡀⠉⢳⡤⢬⡇⠀⠀⠀⠀⠀⠀⠀⣴⣋⣀⣀⣀⣀⣀⣉⣷⣦⠿⠀⠀⠀⠀⠀⢹⣆⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⣿⡰⠏⠀⠀⠀⠀⠀⠀⠀⠙⣆⢸⡷⠶⠤⡷⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠉⠉⠁⠈⠉⠀⠀⠀⠀⠀⠀⠀⠀⠻⡿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀
// ⠀⠀⠀⠀⠀⠀⠘⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣛⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
