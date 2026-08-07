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
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

// gridStep returns the snapping-grid cell size in canvas pixels.
func (this *ZoneEditorDialog) gridStep() float64 {
	return this.zoneHandler.GetZoneEditorGridStep(this.radius)
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
	paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.GridLine, clip.Outline{Path: path.End()}.Op())
}

// drawSnapGuides draws thin green lines across the canvas where the dragged
// zone is currently holding onto another zone's edge/center extension. Only
// visible while a zone is being moved.
func (this *ZoneEditorDialog) drawSnapGuides(gtx layout.Context) {
	if this.zoneDragName == "" || !this.zoneDragMoved {
		return
	}
	if this.snapGuideXActive {
		guide := int(math.Round(this.snapGuideX))
		line := clip.Rect{Min: image.Pt(guide, 0), Max: image.Pt(guide+1, this.side)}
		paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.SnapGuide, line.Op())
	}
	if this.snapGuideYActive {
		guide := int(math.Round(this.snapGuideY))
		line := clip.Rect{Min: image.Pt(0, guide), Max: image.Pt(this.side, guide+1)}
		paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.SnapGuide, line.Op())
	}
}

// snapDraggedPosition nudges the dragged zone's center so that its edges or
// center "hold on" to nearby guides, and records the alignment lines it caught
// so the overlay can draw them.
func (this *ZoneEditorDialog) snapDraggedPosition(pos image.Point) image.Point {
	this.snapGuideXActive = false
	this.snapGuideYActive = false
	if !this.snapBool.Value {
		return pos
	}
	result := this.zoneHandler.SnapZoneEditorPosition(dtos.ZoneEditorSnapRequestDto{
		Position:    pos,
		Positions:   this.positions,
		ZoneRadius:  this.radius,
		DraggedZone: this.zoneDragName,
	})
	if result.HasGuideX {
		this.snapGuideX, this.snapGuideXActive = result.GuideX, true
	}
	if result.HasGuideY {
		this.snapGuideY, this.snapGuideYActive = result.GuideY, true
	}
	return result.Position
}
