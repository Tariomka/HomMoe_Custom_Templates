package models

import (
	"math"
	"slices"
)

type Vector2 struct {
	X, Y float64
}

func NewPosition(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
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
	for i, label := range orderedLabels {
		tier := 0
		if !slices.Contains(playerLabels, label) {
			tier = neutralZonePlans.GetTier(label)
		}
		byTier[tier] = append(byTier[tier], i)
	}

	positions := newPositions(count)
	for tier, indices := range byTier {
		radius := tierRadius(tier)
		nn := len(indices)
		offset := float64(tier) * math.Pi / math.Max(1, float64(nn))
		for i, idx := range indices {
			angle := 2*math.Pi*float64(i)/float64(nn) + offset
			jitter := float64(i%3-1) * 0.008
			positions[idx] = Vector2{
				X: math.Max(0.05, math.Min(0.95, 0.5+math.Cos(angle+jitter)*radius)),
				Y: math.Max(0.05, math.Min(0.95, 0.5+math.Sin(angle+jitter)*radius)),
			}
		}
	}
	return positions
}

func (this *Positions) Add(position Vector2) { *this = append(*this, position) }

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

// CreateDelaunayTriangulation creates a Delaunay triangulation of the positions,
// which divides the positions into a mesh of contiguous, non-overlapping triangles,
// and returns the edges as pairs of indexes. Delaunay triangulation ensures that no point
// falls inside the circumscribed circle (the circle that touches all three vertices) of any triangle.
func (this *Positions) CreateDelaunayTriangulation() [][2]int {
	count := len(*this)
	if count <= 1 {
		return nil
	}
	if count == 2 {
		return [][2]int{{0, 1}}
	}
	minX, minY := (*this)[0].X, (*this)[0].Y
	maxX, maxY := minX, minY
	for _, p := range (*this)[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	dx, dy := maxX-minX, maxY-minY
	delta := math.Max(dx, dy) * 10
	superPts := make([][2]float64, count+3)
	for i, p := range *this {
		superPts[i] = [2]float64{p.X, p.Y}
	}
	superPts[count] = [2]float64{minX - delta, minY - delta*3}
	superPts[count+1] = [2]float64{minX + delta*3, minY - delta}
	superPts[count+2] = [2]float64{minX, minY + delta*3}

	type tri struct{ i0, i1, i2 int }
	triangles := []tri{{count, count + 1, count + 2}}

	for p := 0; p < count; p++ {
		px, py := superPts[p][0], superPts[p][1]
		var bad []tri
		for _, t := range triangles {
			if inCircumscribedCircle(superPts, t.i0, t.i1, t.i2, px, py) {
				bad = append(bad, t)
			}
		}
		type edge struct{ a, b int }
		var boundary []edge
		for _, t := range bad {
			edges := [3]edge{{t.i0, t.i1}, {t.i1, t.i2}, {t.i2, t.i0}}
			for _, e := range edges {
				shared := false
				for _, o := range bad {
					if o == t {
						continue
					}
					if (o.i0 == e.a && o.i1 == e.b) || (o.i1 == e.a && o.i0 == e.b) ||
						(o.i1 == e.a && o.i2 == e.b) || (o.i2 == e.a && o.i1 == e.b) ||
						(o.i2 == e.a && o.i0 == e.b) || (o.i0 == e.a && o.i2 == e.b) {
						shared = true
						break
					}
				}
				if !shared {
					boundary = append(boundary, e)
				}
			}
		}
		badSet := map[tri]bool{}
		for _, t := range bad {
			badSet[t] = true
		}
		var newTris []tri
		for _, t := range triangles {
			if !badSet[t] {
				newTris = append(newTris, t)
			}
		}
		for _, e := range boundary {
			newTris = append(newTris, tri{e.a, e.b, p})
		}
		triangles = newTris
	}

	var realTris []tri
	for _, t := range triangles {
		if t.i0 < count && t.i1 < count && t.i2 < count {
			realTris = append(realTris, t)
		}
	}
	edgeSet := map[[2]int]bool{}
	for _, t := range realTris {
		addEdge := func(a, b int) {
			if a > b {
				a, b = b, a
			}
			edgeSet[[2]int{a, b}] = true
		}
		addEdge(t.i0, t.i1)
		addEdge(t.i1, t.i2)
		addEdge(t.i2, t.i0)
	}
	result := make([][2]int, 0, len(edgeSet))
	for e := range edgeSet {
		result = append(result, e)
	}
	return result

}

// inCircumscribedCircle reports whether (px,py) lies strictly inside the
// circumscribed circle of the triangle (pts[i0], pts[i1], pts[i2]).
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
//	⠀⡇⠀⠀⠀⠀⠀⢀⡜⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢣⡀⠀⠀⠀⠀⠀⢸⠁
//	⠀⡇⠀⠀⠀⠀⢀⡞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢳⡀⠀⠀⠀⠀⢸
//	⠀⢹⠀⠀⠀⢠⠎⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠱⡄⠀⠀⠀⡏
//	⠀⠈⢧⠀⢠⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⡄⠀⡼⠁
//	⠀⠀⠈⢶⣃⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣘⡶⠁
//	⠀⠀⠀⠈⢣⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡜⠁
//	⠀⠀⠀⠀⠀⠙⢦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡴⠋
//	⠀⠀⠀⠀⠀⠀⠀⠙⠦⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⠴⠋
//	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠢⢤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⡤⠔⠋⠁
//	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠙⠒⠒⠒⠒⠖⠒⠒⠒⠋⠉⠁
//	```
func inCircumscribedCircle(pts [][2]float64, i0, i1, i2 int, px, py float64) bool {
	ax, ay := pts[i0][0]-px, pts[i0][1]-py
	bx, by := pts[i1][0]-px, pts[i1][1]-py
	cx, cy := pts[i2][0]-px, pts[i2][1]-py
	det := ax*(by*(cx*cx+cy*cy)-cy*(bx*bx+by*by)) -
		ay*(bx*(cx*cx+cy*cy)-cx*(bx*bx+by*by)) +
		(ax*ax+ay*ay)*(bx*cy-by*cx)
	return det > 0
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
