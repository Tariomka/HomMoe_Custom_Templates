package dialogs

import (
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"gioui.org/f32"
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
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
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
	this.side = side
	this.handlePointer(gtx)

	// Every mutator raises geometryDirty, so an idle dialog skips the full
	// BuildPreviewLayout + edge-grouping pass and redraws cached geometry.
	if this.geometryDirty || this.geometrySide != side {
		this.recomputeGeometry(side)
		this.geometryDirty = false
		this.geometrySide = side
	}

	if len(this.positions) == 0 {
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
		pos := image.Pt(int(pe.Position.X), int(pe.Position.Y))
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

func (this *ZoneEditorDialog) onPress(pos image.Point, pe pointer.Event) {
	node := this.hitTestNode(pos)
	if pe.Buttons&pointer.ButtonSecondary != 0 {
		if edge := this.hitTestEdge(pos); edge != nil {
			this.deleteConnection(edge)
		}
		return
	}
	if this.addMode {
		if node != "" {
			this.pendingFrom = node
			this.dragging = true
			this.dragPos = pos
		} else {
			this.addMode = false
			this.pendingFrom = ""
			this.dragging = false
		}
		return
	}
	if this.addZoneMode {
		if node == "" {
			this.addZoneAt(pos)
		} else {
			this.addZoneMode = false
		}
		return
	}
	if node != "" {
		this.selectZone(node)
		this.zoneDragName = node
		this.zoneDragMoved = false
		this.pressPos = pos
		return
	}
	if edge := this.hitTestEdge(pos); edge != nil {
		this.selected = edge
		this.selectedZone = ""
		this.syncedFor = nil
	} else {
		this.selected = nil
		this.selectedZone = ""
	}
}

func (this *ZoneEditorDialog) onRelease(pos image.Point) {
	this.zoneDragName = ""
	this.zoneDragMoved = false
	if !this.dragging {
		return
	}
	this.dragging = false
	from := this.pendingFrom
	this.pendingFrom = ""
	if !this.addMode {
		return
	}
	target := this.hitTestNode(pos)
	if target != "" && from != "" && target != from {
		// Stay in add mode so several connections can be chained without
		// re-clicking the toolbar button.
		this.addConnection(from, target)
	}
}

func (this *ZoneEditorDialog) hitTestNode(pos image.Point) string {
	best := ""
	bestDistance := math.MaxFloat64
	reach := float64(this.radius)
	for name, center := range this.positions {
		distance := math.Hypot(float64(pos.X-center.X), float64(pos.Y-center.Y))
		if distance <= reach && distance < bestDistance {
			bestDistance = distance
			best = name
		}
	}
	return best
}

func (this *ZoneEditorDialog) hitTestEdge(pos image.Point) *entities.Connection {
	var best *entities.Connection
	bestDistance := 9.0
	for i := range this.edges {
		edge := this.edges[i]
		for step := range 21 {
			t := float64(step) / 20.0
			bezierPoint := helpers.GetVectorOnQuadraticBezierCurve(
				data.NewVec2(float64(edge.startPoint.X), float64(edge.startPoint.Y)),
				data.NewVec2(float64(edge.controlPoint.X), float64(edge.controlPoint.Y)),
				data.NewVec2(float64(edge.endPoint.X), float64(edge.endPoint.Y)),
				t)
			distance := math.Hypot(float64(pos.X)-bezierPoint.X, float64(pos.Y)-bezierPoint.Y)
			if distance < bestDistance {
				bestDistance = distance
				best = edge.connection
			}
		}
	}
	return best
}

// recomputeGeometry rebuilds node positions (via BuildPreviewLayout, identical to
// the preview tab) and curved-edge control points, spreading parallel edges and
// bulging around intermediate nodes.
func (this *ZoneEditorDialog) recomputeGeometry(side int) {
	this.side = side
	response, err := this.previewHandler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{
		Zones:       this.zones,
		Connections: derefConnections(this.working),
		Topology:    this.topology,
		CanvasSide:  float64(side),
	})
	if err != nil {
		this.positions = map[string]image.Point{}
		this.previewZones = nil
		this.radius = 0
		return
	}
	layoutData := response.Layout
	this.positions = layoutData.Positions
	this.previewZones = layoutData.Zones
	this.radius = layoutData.ZoneRadius

	order, groups := this.groupConnectionsByPair()

	this.edges = this.edges[:0]
	const bulgeGap = 18.0
	for _, key := range order {
		connections := groups[key]
		count := len(connections)
		for index, connection := range connections {
			p0, ok0 := this.positions[connection.From]
			p1, ok1 := this.positions[connection.To]
			if !ok0 || !ok1 {
				continue
			}
			canonicalA, canonicalB := p0, p1
			if connection.From > connection.To {
				canonicalA, canonicalB = canonicalB, canonicalA
			}
			dx := float64(canonicalB.X - canonicalA.X)
			dy := float64(canonicalB.Y - canonicalA.Y)
			distance := math.Hypot(dx, dy)
			if distance < 1 {
				distance = 1
			}
			normalX := dy / distance
			normalY := -dx / distance
			spread := (float64(index) - float64(count-1)/2.0) * bulgeGap
			bulge := spread + this.obstacleBulge(canonicalA, canonicalB, normalX, normalY)
			midX := float64(p0.X+p1.X) / 2.0
			midY := float64(p0.Y+p1.Y) / 2.0
			ctrlX := midX + 2.0*bulge*normalX
			ctrlY := midY + 2.0*bulge*normalY
			labelX := 0.25*float64(p0.X) + 0.5*ctrlX + 0.25*float64(p1.X)
			labelY := 0.25*float64(p0.Y) + 0.5*ctrlY + 0.25*float64(p1.Y)
			this.edges = append(this.edges, connectionEdgeGeometry{
				connection:   connection,
				startPoint:   f32.Pt(float32(p0.X), float32(p0.Y)),
				endPoint:     f32.Pt(float32(p1.X), float32(p1.Y)),
				controlPoint: f32.Pt(float32(ctrlX), float32(ctrlY)),
				midPoint:     image.Pt(int(labelX), int(labelY)),
			})
		}
	}
}

// groupConnectionsByPair buckets the working connections by unordered endpoint
// pair, preserving first-seen order so parallel edges spread deterministically
// from frame to frame.
func (this *ZoneEditorDialog) groupConnectionsByPair() ([]connectionPairKey, map[connectionPairKey][]*entities.Connection) {
	groups := make(map[connectionPairKey][]*entities.Connection)
	order := make([]connectionPairKey, 0)
	for _, connection := range this.working {
		a, b := connection.From, connection.To
		if a > b {
			a, b = b, a
		}
		key := connectionPairKey{a, b}
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], connection)
	}
	return order, groups
}

// obstacleBulge returns a perpendicular push so the curve bends clear of any zone
// node that lies close to the straight chord between its two endpoints.
func (this *ZoneEditorDialog) obstacleBulge(a, b image.Point, normalX, normalY float64) float64 {
	clearance := float64(this.radius) + 8.0
	ax, ay := float64(a.X), float64(a.Y)
	segX := float64(b.X - a.X)
	segY := float64(b.Y - a.Y)
	segLength2 := segX*segX + segY*segY
	if segLength2 < 1 {
		return 0
	}
	best := 0.0
	bestMagnitude := 0.0
	for _, center := range this.positions {
		px, py := float64(center.X), float64(center.Y)
		t := ((px-ax)*segX + (py-ay)*segY) / segLength2
		if t <= 0.08 || t >= 0.92 {
			continue
		}
		closestX := ax + t*segX
		closestY := ay + t*segY
		perpendicular := math.Hypot(px-closestX, py-closestY)
		if perpendicular >= clearance {
			continue
		}
		side := (px-closestX)*normalX + (py-closestY)*normalY
		need := (clearance - perpendicular) + 6.0
		signed := need
		if side >= 0 {
			signed = -need
		}
		if math.Abs(signed) > bestMagnitude {
			bestMagnitude = math.Abs(signed)
			best = signed
		}
	}
	return best
}

func (this *ZoneEditorDialog) drawEdges(gtx layout.Context, theme *material.Theme) {
	for i := range this.edges {
		edge := this.edges[i]
		lineColor := themes.ColorsPreview.DirectLine
		width := float32(gtx.Dp(unit.Dp(2)))
		if strings.EqualFold(edge.connection.ConnectionType, "Portal") {
			lineColor = themes.ColorsPreview.PortalLine
			width = float32(gtx.Dp(unit.Dp(1.6)))
		}
		if edge.connection == this.selected {
			lineColor = themes.ColorsZoneEditor.EdgeSelected
			width = float32(gtx.Dp(unit.Dp(3)))
		}
		var path clip.Path
		path.Begin(gtx.Ops)
		path.MoveTo(edge.startPoint)
		path.QuadTo(edge.controlPoint, edge.endPoint)
		paint.FillShape(gtx.Ops, lineColor, clip.Stroke{Path: path.End(), Width: width}.Op())

		if edge.connection.IsUserAdded {
			marker := gtx.Dp(unit.Dp(3))
			dot := image.Rect(
				edge.midPoint.X-marker,
				edge.midPoint.Y-marker,
				edge.midPoint.X+marker,
				edge.midPoint.Y+marker,
			)
			paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.UserAddedDot, clip.UniformRRect(dot, marker).Op(gtx.Ops))
		}
		drawCanvasText(
			gtx,
			theme,
			image.Pt(edge.midPoint.X, edge.midPoint.Y-gtx.Dp(unit.Dp(9))),
			strconv.Itoa(edge.connection.GuardValue),
			9,
			themes.ColorsZoneEditor.GuardLabel,
		)
	}
}

func (this *ZoneEditorDialog) drawRubberBand(gtx layout.Context) {
	if !this.addMode || !this.dragging || this.pendingFrom == "" {
		return
	}
	center, ok := this.positions[this.pendingFrom]
	if !ok {
		return
	}
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(float32(center.X), float32(center.Y)))
	path.LineTo(f32.Pt(float32(this.dragPos.X), float32(this.dragPos.Y)))
	paint.FillShape(
		gtx.Ops,
		themes.ColorsZoneEditor.EdgeSelected,
		clip.Stroke{Path: path.End(), Width: float32(gtx.Dp(unit.Dp(2)))}.Op(),
	)
}

func (this *ZoneEditorDialog) drawNodes(gtx layout.Context, theme *material.Theme) {
	for _, zone := range this.previewZones {
		if zone.Type == preview.ZoneTypePlayer {
			continue
		}
		utils.DrawPreviewZone(gtx, theme, zone, this.radius)
	}
	for _, zone := range this.previewZones {
		if zone.Type != preview.ZoneTypePlayer {
			continue
		}
		utils.DrawPreviewZone(gtx, theme, zone, this.radius)
	}
	if this.addMode && this.pendingFrom != "" {
		if center, ok := this.positions[this.pendingFrom]; ok {
			reach := this.radius + 4
			rect := image.Rect(center.X-reach, center.Y-reach, center.X+reach, center.Y+reach)
			paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.EdgeSelected, clip.Stroke{
				Path:  clip.UniformRRect(rect, reach).Path(gtx.Ops),
				Width: float32(gtx.Dp(unit.Dp(2))),
			}.Op())
		}
	}
	if !this.addMode && this.selectedZone != "" {
		if center, ok := this.positions[this.selectedZone]; ok {
			reach := this.radius + 4
			rect := image.Rect(center.X-reach, center.Y-reach, center.X+reach, center.Y+reach)
			paint.FillShape(gtx.Ops, themes.ColorsZoneEditor.EdgeSelected, clip.Stroke{
				Path:  clip.UniformRRect(rect, reach).Path(gtx.Ops),
				Width: float32(gtx.Dp(unit.Dp(2))),
			}.Op())
		}
	}
}

// moveDraggedZone updates the dragged zone's manual position. The drag only
// starts once the pointer moved a few pixels from the press point, so plain
// clicks still just select.
func (this *ZoneEditorDialog) moveDraggedZone(pos image.Point) {
	if this.addMode || this.zoneDragName == "" || this.side <= 0 {
		return
	}
	if !this.zoneDragMoved {
		dx := float64(pos.X - this.pressPos.X)
		dy := float64(pos.Y - this.pressPos.Y)
		if math.Hypot(dx, dy) < 6 {
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
	x := math.Min(math.Max(float64(pos.X)/float64(this.side), 0.04), 0.96)
	y := math.Min(math.Max(float64(pos.Y)/float64(this.side), 0.04), 0.96)
	zone.ManualPosition = &[2]float64{x, y}
	this.geometryDirty = true
}

func drawCanvasText(
	gtx layout.Context,
	theme *material.Theme,
	center image.Point,
	text string,
	sizeSp int,
	textColor color.NRGBA,
) {
	macro := op.Record(gtx.Ops)
	local := gtx
	local.Constraints.Min = image.Point{}
	local.Constraints.Max = image.Pt(1<<14, 1<<14)
	label := material.Label(theme, unit.Sp(float32(sizeSp)), text)
	label.Color = textColor
	label.Font = font.Font{Weight: font.Medium}
	dims := label.Layout(local)
	call := macro.Stop()
	offset := op.Offset(image.Pt(center.X-dims.Size.X/2, center.Y-dims.Size.Y/2)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	offset.Pop()
}
