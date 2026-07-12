package dialogs

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

// Snapping tunables - adjust the feel of the snap toggle here.
const (
	// gridCellsPerZoneDiameter sets the snap-grid density: the distance
	// between adjacent grid lines is (zone diameter) / gridCellsPerZoneDiameter.
	gridCellsPerZoneDiameter = 7.0
	// gridSnapThresholdPx is the "light" hold distance (canvas px) within
	// which a dragged zone's edges/centre stick to a grid line.
	gridSnapThresholdPx = 4.0
	// zoneSnapThresholdPx is the "heavier" hold distance (canvas px) within
	// which a dragged zone's edges/centre stick to the horizontal or vertical
	// extension of another zone's edge or centre.
	zoneSnapThresholdPx = 9.0
)

// gridStep returns the snapping-grid cell size in canvas pixels.
func (this *ZoneEditorDialog) gridStep() float64 {
	return float64(this.radius) * 2.0 / gridCellsPerZoneDiameter
}

// drawSnapGrid paints faint dots at the snapping-grid intersections behind
// edges and nodes while the snap toggle is on.
func (this *ZoneEditorDialog) drawSnapGrid(gtx layout.Context) {
	if !this.snapBool.Value {
		return
	}
	step := this.gridStep()
	if step < 3 {
		return
	}
	half := float32(gtx.Dp(unit.Dp(1)))
	// All dots go into a single path so the grid costs one fill op per frame.
	var path clip.Path
	path.Begin(gtx.Ops)
	for x := step; x < float64(this.side); x += step {
		for y := step; y < float64(this.side); y += step {
			fx, fy := float32(x), float32(y)
			path.MoveTo(f32.Pt(fx-half, fy-half))
			path.LineTo(f32.Pt(fx+half, fy-half))
			path.LineTo(f32.Pt(fx+half, fy+half))
			path.LineTo(f32.Pt(fx-half, fy+half))
			path.Close()
		}
	}
	paint.FillShape(gtx.Ops, themes.ColorEditorGridLine, clip.Outline{Path: path.End()}.Op())
}

// drawSnapGuides draws thin green lines across the canvas where the dragged
// zone is currently holding onto another zone's edge/centre extension. Only
// visible while a zone is being moved.
func (this *ZoneEditorDialog) drawSnapGuides(gtx layout.Context) {
	if this.zoneDragName == "" || !this.zoneDragMoved {
		return
	}
	if this.snapGuideXActive {
		guide := int(math.Round(this.snapGuideX))
		line := clip.Rect{Min: image.Pt(guide, 0), Max: image.Pt(guide+1, this.side)}
		paint.FillShape(gtx.Ops, themes.ColorEditorSnapGuide, line.Op())
	}
	if this.snapGuideYActive {
		guide := int(math.Round(this.snapGuideY))
		line := clip.Rect{Min: image.Pt(0, guide), Max: image.Pt(this.side, guide+1)}
		paint.FillShape(gtx.Ops, themes.ColorEditorSnapGuide, line.Op())
	}
}

// snapDraggedPosition nudges the dragged zone's centre so that its edges or
// centre "hold on" to nearby guides: heavier onto other zones' edge/centre
// extension lines, lighter onto the background grid. It never pulls from afar -
// only positions already within the threshold stick.
func (this *ZoneEditorDialog) snapDraggedPosition(pos image.Point) image.Point {
	this.snapGuideXActive = false
	this.snapGuideYActive = false
	if !this.snapBool.Value || this.radius <= 0 {
		return pos
	}
	radius := float64(this.radius)
	// The dragged zone's own snap points on each axis: leading edge, centre,
	// trailing edge.
	offsets := [3]float64{-radius, 0, radius}
	guidesX, guidesY := this.otherZoneGuides(radius)
	x, guideX, hitX := snapAxis(float64(pos.X), offsets, guidesX, this.gridStep())
	y, guideY, hitY := snapAxis(float64(pos.Y), offsets, guidesY, this.gridStep())
	if hitX {
		this.snapGuideX, this.snapGuideXActive = guideX, true
	}
	if hitY {
		this.snapGuideY, this.snapGuideYActive = guideY, true
	}
	return image.Pt(int(math.Round(x)), int(math.Round(y)))
}

// otherZoneGuides collects the horizontal and vertical guide coordinates
// (edge / centre / edge) of every zone except the dragged one.
func (this *ZoneEditorDialog) otherZoneGuides(radius float64) (guidesX, guidesY []float64) {
	for name, center := range this.positions {
		if name == this.zoneDragName {
			continue
		}
		cx, cy := float64(center.X), float64(center.Y)
		guidesX = append(guidesX, cx-radius, cx, cx+radius)
		guidesY = append(guidesY, cy-radius, cy, cy+radius)
	}
	return guidesX, guidesY
}

// snapAxis snaps a single axis value. Zone-alignment guides win over the grid;
// within each class the smallest correction wins. When a zone guide is hit its
// coordinate is returned so the caller can draw an alignment indicator.
func snapAxis(
	value float64,
	offsets [3]float64,
	guides []float64,
	gridStep float64,
) (snapped float64, guide float64, zoneGuideHit bool) {
	best := math.MaxFloat64
	bestGuide := 0.0
	for _, offset := range offsets {
		point := value + offset
		for _, g := range guides {
			if delta := g - point; math.Abs(delta) < math.Abs(best) {
				best = delta
				bestGuide = g
			}
		}
	}
	if math.Abs(best) <= zoneSnapThresholdPx {
		return value + best, bestGuide, true
	}
	if gridStep > 0 {
		best = math.MaxFloat64
		for _, offset := range offsets {
			point := value + offset
			if delta := math.Round(point/gridStep)*gridStep - point; math.Abs(delta) < math.Abs(best) {
				best = delta
			}
		}
		if math.Abs(best) <= gridSnapThresholdPx {
			return value + best, 0, false
		}
	}
	return value, 0, false
}
