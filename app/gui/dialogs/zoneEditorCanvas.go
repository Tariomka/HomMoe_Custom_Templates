package dialogs

import (
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

// layoutCanvas draws the node/edge canvas and processes pointer interaction. All
// coordinates are square-local because the centring offset is pushed first and
// the pointer area is registered within that transform.
func (this *ZoneEditorDialog) layoutCanvas(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	maxX := gtx.Constraints.Max.X
	maxY := gtx.Constraints.Max.Y
	outer := image.Pt(maxX, maxY)
	side := max(min(maxX, maxY), 80)
	canvasSize := image.Pt(side, side)
	offsetX := (maxX - side) / 2
	offsetY := (maxY - side) / 2
	defer op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops).Pop()

	paint.FillShape(gtx.Ops, themes.ColorsPreview.Background, clip.Rect(image.Rectangle{Max: canvasSize}).Op())
	frameRadius := gtx.Dp(unit.Dp(6))
	frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
	paint.FillShape(gtx.Ops, themes.ColorsPreview.Frame, clip.Stroke{
		Path:  clip.UniformRRect(frame, frameRadius).Path(gtx.Ops),
		Width: 2,
	}.Op())

	if len(this.zones) == 0 {
		return widgets.NewCenteredMessageWidget(theme,
			"No zones to edit - generate a template first.", canvasSize, outer)(gtx)
	}

	area := clip.Rect{Max: canvasSize}.Push(gtx.Ops)
	event.Op(gtx.Ops, &this.canvasTag)
	area.Pop()
	// Handle pointer input BEFORE recomputing geometry so edits (new
	// connections/zones, drags) are reflected in this very frame. Hit testing
	// uses the previous frame's geometry, which is exactly what is on screen.
	sideChanged := this.side != side
	this.side = side
	this.handlePointer(gtx)

	// Every mutator raises geometryDirty, so an idle dialog skips the full
	// BuildPreviewLayout + edge-grouping pass and redraws cached geometry.
	if this.geometryDirty || sideChanged {
		this.recomputeGeometry(side)
		this.geometryDirty = false
	}

	if len(this.geometry.Positions) == 0 {
		return layout.Dimensions{Size: outer}
	}

	this.drawSnapGrid(gtx)
	this.drawEdges(gtx, theme)
	this.drawRubberBand(gtx)
	this.drawNodes(gtx, theme)
	this.drawSnapGuides(gtx)
	return layout.Dimensions{Size: outer}
}

func (this *ZoneEditorDialog) handlePointer(gtx layout.Context) {
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: &this.canvasTag,
			Kinds:  pointer.Press | pointer.Drag | pointer.Release,
		})
		if !ok {
			break
		}

		pe, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		pos := utils.ToVec2(pe.Position)
		switch pe.Kind {
		case pointer.Press:
			this.onPress(pos, pe)
		case pointer.Drag:
			this.dragPos = pos
			this.moveDraggedZone(pos)
		case pointer.Release:
			this.onRelease(pos)
		case pointer.Cancel, pointer.Enter, pointer.Leave, pointer.Move, pointer.Scroll: // noop
		}
	}
}

func (this *ZoneEditorDialog) onPress(pos models.Position, pe pointer.Event) {
	node := this.hitTestNode(pos)
	if pe.Buttons&pointer.ButtonSecondary != 0 {
		if edge := this.hitTestEdge(pos); edge != nil {
			this.deleteConnection(edge)
		}
		return
	}

	if this.addMode {
		if node != "" {
			this.beginConnectionDrag(node, pos)
		} else {
			this.exitAddModes()
		}
		return
	}

	if this.addZoneMode {
		if node == "" {
			this.addZoneAt(pos)
		} else {
			this.exitAddModes()
		}
		return
	}

	if node != "" {
		this.selectZone(node)
		this.beginZoneDrag(node, pos)
		return
	}

	if edge := this.hitTestEdge(pos); edge != nil {
		this.selectConnection(edge)
		this.syncedFor = nil
	} else {
		this.clearSelection()
	}
}

func (this *ZoneEditorDialog) onRelease(pos models.Position) {
	this.endZoneDrag()
	from := this.finishConnectionDrag()
	if !this.addMode || from == "" {
		return
	}

	target := this.hitTestNode(pos)
	if target != "" && target != from {
		// Stay in add mode so several connections can be chained without
		// re-clicking the toolbar button.
		this.addConnection(from, target)
	}
}

func (this *ZoneEditorDialog) hitTestNode(pos models.Position) string {
	return this.zoneHandler.HitTestZoneEditorNode(dtos.ZoneEditorHitTestRequestDto{
		Position:   pos,
		Positions:  this.geometry.Positions,
		ZoneRadius: this.geometry.ZoneRadius,
	})
}

func (this *ZoneEditorDialog) hitTestEdge(pos models.Position) *entities.Connection {
	edgeIndex := this.zoneHandler.HitTestZoneEditorEdge(pos, this.geometry.Edges)
	if edgeIndex < 0 {
		return nil
	}

	return this.edgeConnection(this.geometry.Edges[edgeIndex])
}

// edgeConnection resolves the working connection an edge was laid out for, or
// nil when the cached geometry no longer lines up with the connection list.
func (this *ZoneEditorDialog) edgeConnection(edge models.ZoneEditorEdge) *entities.Connection {
	if edge.ConnectionIndex < 0 || edge.ConnectionIndex >= len(this.working) {
		return nil
	}

	return this.working[edge.ConnectionIndex]
}

// recomputeGeometry rebuilds node positions and curved-edge control points via
// the geometry service, which places nodes exactly as the preview tab does.
func (this *ZoneEditorDialog) recomputeGeometry(side int) {
	this.side = side
	this.geometry = this.zoneHandler.BuildZoneEditorGeometry(dtos.ZoneEditorGeometryRequestDto{
		Zones:       this.zones,
		Connections: derefConnections(this.working),
		Topology:    this.topology,
		CanvasSide:  side,
	})
}

func (this *ZoneEditorDialog) drawEdges(gtx layout.Context, theme *material.Theme) {
	for i := range this.geometry.Edges {
		edge := this.geometry.Edges[i]
		connection := this.edgeConnection(edge)
		if connection == nil {
			continue
		}

		lineColor := themes.ColorsPreview.DirectLine
		width := float32(gtx.Dp(unit.Dp(2)))
		if strings.EqualFold(connection.ConnectionType, "Portal") {
			lineColor = themes.ColorsPreview.PortalLine
			width = float32(gtx.Dp(unit.Dp(1.6)))
		}
		if connection == this.selected {
			lineColor = themes.ColorsZoneEditor.EdgeSelected
			width = float32(gtx.Dp(unit.Dp(3)))
		}
		var path clip.Path
		path.Begin(gtx.Ops)
		path.MoveTo(utils.ToF32Point(edge.StartPoint))
		path.QuadTo(utils.ToF32Point(edge.ControlPoint), utils.ToF32Point(edge.EndPoint))
		paint.FillShape(gtx.Ops, lineColor, clip.Stroke{Path: path.End(), Width: width}.Op())

		if connection.IsUserAdded {
			marker := gtx.Dp(unit.Dp(3))
			mid := edge.MidPoint.ToPointRounded()
			dot := image.Rect(
				mid.X-marker, mid.Y-marker,
				mid.X+marker, mid.Y+marker)
			paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.UserAddedDot, clip.UniformRRect(dot, marker).Op(gtx.Ops))
		}
		drawCanvasText(
			gtx, theme, edge.MidPoint.Subtract(data.NewVec2(0, float64(gtx.Dp(unit.Dp(9))))),
			strconv.Itoa(connection.GuardValue), 9, themes.ColorsZoneEditor.GuardLabel)
	}
}

func (this *ZoneEditorDialog) drawRubberBand(gtx layout.Context) {
	if !this.addMode || !this.dragging || this.pendingFrom == "" {
		return
	}

	center, ok := this.geometry.Positions[this.pendingFrom]
	if !ok {
		return
	}

	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(utils.ToF32Point(center))
	path.LineTo(utils.ToF32Point(this.dragPos))
	paint.FillShape(
		gtx.Ops, themes.ColorsZoneEditor.EdgeSelected,
		clip.Stroke{Path: path.End(), Width: float32(gtx.Dp(unit.Dp(2)))}.Op())
}

func (this *ZoneEditorDialog) drawNodes(gtx layout.Context, theme *material.Theme) {
	for _, zone := range this.geometry.Zones {
		if zone.Type == preview.ZoneTypePlayer {
			continue
		}

		utils.DrawPreviewZone(gtx, theme, zone, this.geometry.ZoneRadius)
	}
	for _, zone := range this.geometry.Zones {
		if zone.Type != preview.ZoneTypePlayer {
			continue
		}

		utils.DrawPreviewZone(gtx, theme, zone, this.geometry.ZoneRadius)
	}
	if this.addMode && this.pendingFrom != "" {
		if center, ok := this.geometry.Positions[this.pendingFrom]; ok {
			this.drawSelectionRing(gtx, center)
		}
	}
	if !this.addMode && this.selectedZone != "" {
		if center, ok := this.geometry.Positions[this.selectedZone]; ok {
			this.drawSelectionRing(gtx, center)
		}
	}
}

// drawSelectionRing outlines a node that is selected or is the source of a
// pending connection. The float centre is snapped here because Gio's rounded
// rectangles are defined on the integer pixel grid.
func (this *ZoneEditorDialog) drawSelectionRing(gtx layout.Context, center models.Position) {
	reach := int(math.Round(this.geometry.ZoneRadius)) + 4
	pixelCenter := center.ToPointRounded()
	rect := image.Rect(
		pixelCenter.X-reach, pixelCenter.Y-reach,
		pixelCenter.X+reach, pixelCenter.Y+reach)
	paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.EdgeSelected, clip.Stroke{
		Path:  clip.UniformRRect(rect, reach).Path(gtx.Ops),
		Width: float32(gtx.Dp(unit.Dp(2))),
	}.Op())
}

// moveDraggedZone updates the dragged zone's manual position. The drag only
// starts once the pointer moved a few pixels from the press point, so plain
// clicks still just select.
func (this *ZoneEditorDialog) moveDraggedZone(pos models.Position) {
	if this.addMode || this.zoneDragName == "" || this.side <= 0 {
		return
	}
	if !this.zoneDragMoved {
		if !this.zoneDragLeftDeadZone(pos) {
			return
		}
		this.ensureManualPositions()
		this.zoneDragMoved = true
	}
	zone := this.zoneByName(this.zoneDragName)
	if zone == nil {
		return
	}
	pos = this.snapDraggedPosition(pos)
	x := math.Min(math.Max(pos.X/float64(this.side), 0.04), 0.96)
	y := math.Min(math.Max(pos.Y/float64(this.side), 0.04), 0.96)
	zone.ManualPosition = &[2]float64{x, y}
	this.geometryDirty = true
}

func drawCanvasText(
	gtx layout.Context,
	theme *material.Theme,
	center models.Position,
	text string,
	sizeSp int,
	textColor color.NRGBA) {
	macro := op.Record(gtx.Ops)
	local := gtx
	local.Constraints.Min = image.Point{}
	local.Constraints.Max = image.Pt(1<<14, 1<<14)
	label := material.Label(theme, unit.Sp(float32(sizeSp)), text)
	label.Color = textColor
	label.Font = font.Font{Weight: font.Medium}
	dims := label.Layout(local)
	call := macro.Stop()
	pixelCenter := center.ToPointRounded()
	offset := op.Offset(image.Pt(pixelCenter.X-dims.Size.X/2, pixelCenter.Y-dims.Size.Y/2)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	offset.Pop()
}
