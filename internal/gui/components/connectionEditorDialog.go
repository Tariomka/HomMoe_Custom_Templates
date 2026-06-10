package components

import (
	"fmt"
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
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

var (
	colorEdgeSelected = color.NRGBA{R: 0xFF, G: 0x8C, B: 0x00, A: 0xFF}
	colorUserAddedDot = color.NRGBA{R: 0xFF, G: 0xD2, B: 0x80, A: 0xFF}
	colorGuardLabel   = color.NRGBA{R: 0xF0, G: 0xE6, B: 0x9A, A: 0xFF}
)

// connEdgeGeom is the per-frame drawn geometry of one connection: a quadratic
// Bézier from p0 to p1 bulged through ctrl, with mid the curve's midpoint used
// for the guard-value label and the user-added marker.
type connEdgeGeom struct {
	conn *template.Connection
	p0   f32.Point
	p1   f32.Point
	ctrl f32.Point
	mid  image.Point
}

// ConnectionEditorDialog is the visual zone-connection editor. It renders zones
// as tier-coloured nodes and connections as curved lines, supports drag-to-create
// and right-click-delete, and exposes a property panel for the selected edge. It
// works on a private copy of the connection list; Apply commits, Cancel/✕ discards.
type ConnectionEditorDialog struct {
	zones       []template.Zone
	playerZones map[string]bool
	topology    config.MapTopology
	working     []*template.Connection
	original    []template.Connection
	onApply     func([]template.Connection)

	// Geometry recomputed every frame from BuildPreviewLayout.
	positions    map[string]image.Point
	previewZones []services.PreviewZone
	radius       int
	edges        []connEdgeGeom

	// Canvas interaction.
	canvasTag   int
	selected    *template.Connection
	addMode     bool
	pendingFrom string
	dragging    bool
	dragPos     image.Point

	// Toolbar / footer.
	addBtn    widget.Clickable
	deleteBtn widget.Clickable
	resetBtn  widget.Clickable
	applyBtn  widget.Clickable
	cancelBtn widget.Clickable

	// Property panel.
	sideScroll        widget.List
	syncedFor         *template.Connection
	typeDropdown      *content.DropdownSelector
	guardZoneDropdown *content.DropdownSelector
	guardDropdown     *content.DropdownSelector
	guardPresetValues []int
	weeklyDropdown    *content.DropdownSelector
	guardValueEdit    widget.Editor
	weeklyEdit        widget.Editor
	matchGroupEdit    widget.Editor
	advancedBool      widget.Bool
	escapeBool        widget.Bool
	simSquadBool      widget.Bool
	sidePropDelete    widget.Clickable
}

// NewConnectionEditorDialog builds an editor over a copy of the given zones and
// connections. onApply receives the edited connection list when the user applies.
func NewConnectionEditorDialog(zones []template.Zone, connections []template.Connection, topology config.MapTopology, onApply func([]template.Connection)) *ConnectionEditorDialog {
	players := make(map[string]bool)
	for _, zone := range zones {
		if strings.HasPrefix(zone.Name, "Spawn-") {
			players[zone.Name] = true
		}
	}

	dialog := &ConnectionEditorDialog{
		zones:             zones,
		playerZones:       players,
		topology:          topology,
		onApply:           onApply,
		typeDropdown:      content.NewDropdownSelector(connection_editor.UserCreatableConnectionTypes()),
		guardZoneDropdown: content.NewDropdownSelector(nil),
		guardDropdown:     content.NewDropdownSelector(nil),
		weeklyDropdown:    content.NewDropdownSelector(connection_editor.WeeklyIncrementLabels),
	}
	for i := range connections {
		working := connections[i]
		dialog.working = append(dialog.working, &working)
		dialog.original = append(dialog.original, connection_editor.CloneConnection(connections[i], false))
	}
	dialog.sideScroll.Axis = layout.Vertical
	dialog.guardValueEdit.SingleLine = true
	dialog.weeklyEdit.SingleLine = true
	dialog.matchGroupEdit.SingleLine = true
	return dialog
}

func (this *ConnectionEditorDialog) Title() string { return "Zone Connection Editor" }

func (this *ConnectionEditorDialog) PreferredSize() (unit.Dp, unit.Dp) {
	return unit.Dp(1000), unit.Dp(720)
}

func (this *ConnectionEditorDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.applyBtn.Clicked(gtx) {
		if this.onApply != nil {
			this.onApply(derefConnections(this.working))
		}
		return layout.Dimensions{Size: gtx.Constraints.Max}, true
	}
	if this.cancelBtn.Clicked(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Max}, true
	}
	if this.addBtn.Clicked(gtx) {
		this.addMode = !this.addMode
		this.pendingFrom = ""
		this.dragging = false
	}
	if this.deleteBtn.Clicked(gtx) && this.selected != nil {
		this.deleteConnection(this.selected)
	}
	if this.resetBtn.Clicked(gtx) {
		this.resetToOriginal()
	}

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(this.layoutToolbar(theme)),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return this.layoutCanvas(gtx, theme)
				}),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(10)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return this.layoutSidePanel(gtx, theme)
				}),
			)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(this.layoutFooter(theme)),
	)
	return dims, false
}

func (this *ConnectionEditorDialog) layoutToolbar(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		addLabel := "Add connection"
		if this.addMode {
			addLabel = "Adding… (click empty to cancel)"
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if this.addMode {
					return widgets.NewGoldButtonWidget(theme, addLabel, &this.addBtn, false)(gtx)
				}
				return widgets.NewButtonWidget(theme, addLabel, &this.addBtn, false)(gtx)
			}),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Delete selected", &this.deleteBtn, this.selected == nil)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Reset to generated", &this.resetBtn, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(12)),
			layout.Flexed(1, this.layoutStatus(theme)),
		)
	}
}

func (this *ConnectionEditorDialog) layoutStatus(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		connections := derefConnections(this.working)
		if connection_editor.ComputeHasErrors(this.zones, connections) {
			label := material.Body2(theme, "⚠ A connection references a missing zone — fix before export.")
			label.Color = themes.ColorError
			label.TextSize = unit.Sp(12)
			label.MaxLines = 2
			return label.Layout(gtx)
		}
		message := fmt.Sprintf("%d connections", len(connections))
		if this.addMode {
			message = "Add mode: press a zone and drag to another to connect."
		} else if isolated := connection_editor.FindIsolatedZones(this.zones, connections); len(isolated) > 0 {
			message = fmt.Sprintf("%d connections · %d isolated zone(s)", len(connections), len(isolated))
		}
		label := material.Body2(theme, message)
		label.Color = themes.ColorTextDim
		label.TextSize = unit.Sp(12)
		label.MaxLines = 2
		return label.Layout(gtx)
	}
}

func (this *ConnectionEditorDialog) layoutFooter(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, 0)}
			}),
			layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelBtn, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Rigid(widgets.NewGoldButtonWidget(theme, "Apply changes", &this.applyBtn, false)),
		)
	}
}

// layoutCanvas draws the node/edge canvas and processes pointer interaction. All
// coordinates are square-local because the centring offset is pushed first and
// the pointer area is registered within that transform.
func (this *ConnectionEditorDialog) layoutCanvas(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	maxX := gtx.Constraints.Max.X
	maxY := gtx.Constraints.Max.Y
	outer := image.Pt(maxX, maxY)
	side := max(min(maxX, maxY), 80)
	canvasSize := image.Pt(side, side)
	offsetX := (maxX - side) / 2
	offsetY := (maxY - side) / 2
	defer op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops).Pop()

	paint.FillShape(gtx.Ops, themes.ColorPreviewBg, clip.Rect(image.Rectangle{Max: canvasSize}).Op())
	frameRadius := gtx.Dp(unit.Dp(6))
	frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
	paint.FillShape(gtx.Ops, themes.ColorPreviewFrame, clip.Stroke{
		Path:  clip.UniformRRect(frame, frameRadius).Path(gtx.Ops),
		Width: 2,
	}.Op())

	if len(this.zones) == 0 {
		return widgets.NewCenteredMessageWidget(theme, "No zones to edit — generate a template first.", canvasSize, outer)(gtx)
	}

	this.recomputeGeometry(side)

	area := clip.Rect{Max: canvasSize}.Push(gtx.Ops)
	event.Op(gtx.Ops, &this.canvasTag)
	area.Pop()
	this.handlePointer(gtx)

	if len(this.positions) == 0 {
		return layout.Dimensions{Size: outer}
	}

	this.drawEdges(gtx, theme)
	this.drawRubberBand(gtx)
	this.drawNodes(gtx, theme)
	return layout.Dimensions{Size: outer}
}

func (this *ConnectionEditorDialog) handlePointer(gtx layout.Context) {
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
		case pointer.Release:
			this.onRelease(pos)
		}
	}
}

func (this *ConnectionEditorDialog) onPress(pos image.Point, pe pointer.Event) {
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
	if edge := this.hitTestEdge(pos); edge != nil {
		this.selected = edge
		this.syncedFor = nil
	} else if node == "" {
		this.selected = nil
	}
}

func (this *ConnectionEditorDialog) onRelease(pos image.Point) {
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
		this.addConnection(from, target)
		this.addMode = false
	}
}

func (this *ConnectionEditorDialog) hitTestNode(pos image.Point) string {
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

func (this *ConnectionEditorDialog) hitTestEdge(pos image.Point) *template.Connection {
	var best *template.Connection
	bestDistance := 9.0
	for i := range this.edges {
		edge := this.edges[i]
		for step := 0; step <= 20; step++ {
			t := float64(step) / 20.0
			mt := 1 - t
			bx := mt*mt*float64(edge.p0.X) + 2*mt*t*float64(edge.ctrl.X) + t*t*float64(edge.p1.X)
			by := mt*mt*float64(edge.p0.Y) + 2*mt*t*float64(edge.ctrl.Y) + t*t*float64(edge.p1.Y)
			distance := math.Hypot(float64(pos.X)-bx, float64(pos.Y)-by)
			if distance < bestDistance {
				bestDistance = distance
				best = edge.conn
			}
		}
	}
	return best
}

// recomputeGeometry rebuilds node positions (via BuildPreviewLayout, identical to
// the preview tab) and curved-edge control points, spreading parallel edges and
// bulging around intermediate nodes.
func (this *ConnectionEditorDialog) recomputeGeometry(side int) {
	mini := &template.RmgTemplateModel{
		Variants: []template.Variant{{
			Zones:       this.zones,
			Connections: derefConnections(this.working),
		}},
	}
	layoutData := services.BuildPreviewLayout(mini, this.topology, float64(side))
	this.positions = layoutData.Positions
	this.previewZones = layoutData.Zones
	this.radius = layoutData.ZoneRadius

	type pairKey struct{ a, b string }
	groups := make(map[pairKey][]*template.Connection)
	order := make([]pairKey, 0)
	for _, connection := range this.working {
		a, b := connection.From, connection.To
		if a > b {
			a, b = b, a
		}
		key := pairKey{a, b}
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], connection)
	}

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
			this.edges = append(this.edges, connEdgeGeom{
				conn: connection,
				p0:   f32.Pt(float32(p0.X), float32(p0.Y)),
				p1:   f32.Pt(float32(p1.X), float32(p1.Y)),
				ctrl: f32.Pt(float32(ctrlX), float32(ctrlY)),
				mid:  image.Pt(int(labelX), int(labelY)),
			})
		}
	}
}

// obstacleBulge returns a perpendicular push so the curve bends clear of any zone
// node that lies close to the straight chord between its two endpoints.
func (this *ConnectionEditorDialog) obstacleBulge(a, b image.Point, normalX, normalY float64) float64 {
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

func (this *ConnectionEditorDialog) drawEdges(gtx layout.Context, theme *material.Theme) {
	for i := range this.edges {
		edge := this.edges[i]
		lineColor := themes.ColorPreviewDirectLine
		width := float32(gtx.Dp(unit.Dp(2)))
		if strings.EqualFold(edge.conn.ConnectionType, "Portal") {
			lineColor = themes.ColorPreviewPortalLine
			width = float32(gtx.Dp(unit.Dp(1.6)))
		}
		if edge.conn == this.selected {
			lineColor = colorEdgeSelected
			width = float32(gtx.Dp(unit.Dp(3)))
		}
		var path clip.Path
		path.Begin(gtx.Ops)
		path.MoveTo(edge.p0)
		path.QuadTo(edge.ctrl, edge.p1)
		paint.FillShape(gtx.Ops, lineColor, clip.Stroke{Path: path.End(), Width: width}.Op())

		if edge.conn.IsUserAdded {
			marker := gtx.Dp(unit.Dp(3))
			dot := image.Rect(edge.mid.X-marker, edge.mid.Y-marker, edge.mid.X+marker, edge.mid.Y+marker)
			paint.FillShape(gtx.Ops, colorUserAddedDot, clip.UniformRRect(dot, marker).Op(gtx.Ops))
		}
		drawCanvasText(gtx, theme, image.Pt(edge.mid.X, edge.mid.Y-gtx.Dp(unit.Dp(9))), strconv.Itoa(edge.conn.GuardValue), 9, colorGuardLabel)
	}
}

func (this *ConnectionEditorDialog) drawRubberBand(gtx layout.Context) {
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
	paint.FillShape(gtx.Ops, colorEdgeSelected, clip.Stroke{Path: path.End(), Width: float32(gtx.Dp(unit.Dp(2)))}.Op())
}

func (this *ConnectionEditorDialog) drawNodes(gtx layout.Context, theme *material.Theme) {
	for _, zone := range this.previewZones {
		if zone.IsPlayer {
			continue
		}
		utils.DrawPreviewZone(gtx, theme, zone, this.radius)
	}
	for _, zone := range this.previewZones {
		if !zone.IsPlayer {
			continue
		}
		utils.DrawPreviewZone(gtx, theme, zone, this.radius)
	}
	if this.addMode && this.pendingFrom != "" {
		if center, ok := this.positions[this.pendingFrom]; ok {
			reach := this.radius + 4
			rect := image.Rect(center.X-reach, center.Y-reach, center.X+reach, center.Y+reach)
			paint.FillShape(gtx.Ops, colorEdgeSelected, clip.Stroke{
				Path:  clip.UniformRRect(rect, reach).Path(gtx.Ops),
				Width: float32(gtx.Dp(unit.Dp(2))),
			}.Op())
		}
	}
}

func (this *ConnectionEditorDialog) layoutSidePanel(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	width := gtx.Dp(unit.Dp(300))
	size := image.Pt(width, gtx.Constraints.Max.Y)
	radius := gtx.Dp(4)
	rect := image.Rectangle{Max: size}
	paint.FillShape(gtx.Ops, themes.ColorPanel, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, themes.ColorBorder, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())

	inner := gtx
	inner.Constraints = layout.Exact(size)
	layout.UniformInset(unit.Dp(10)).Layout(inner, func(gtx layout.Context) layout.Dimensions {
		if this.selected == nil {
			return this.layoutSideHint(gtx, theme)
		}
		if this.syncedFor != this.selected {
			this.syncPropsFromConnection()
			this.syncedFor = this.selected
		}
		rows := this.propertyRows(theme)
		dims := material.List(theme, &this.sideScroll).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
			return rows[index](gtx)
		})
		this.writebackProps(gtx)
		return dims
	})
	return layout.Dimensions{Size: size}
}

func (this *ConnectionEditorDialog) layoutSideHint(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, "No connection selected")
			label.Color = themes.ColorText
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Click a connection line to edit its guard, type and growth.")),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Use “Add connection”, then drag from one zone to another to create a link.")),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Right-click a line to delete it.")),
	)
}

func (this *ConnectionEditorDialog) propertyRows(theme *material.Theme) []layout.Widget {
	connection := this.selected
	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, connection.From+"  →  "+connection.To)
			label.Color = themes.ColorGoldBright
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		},
	}
	if connection.IsUserAdded {
		rows = append(rows, widgets.NewDimmedLabelWidget(theme, "User-added connection"))
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(6),
		widgets.NewLabeledRowWidget(theme, "Type", 110, func(gtx layout.Context) layout.Dimensions {
			return this.typeDropdown.Layout(gtx, theme)
		}),
		widgets.NewLabeledRowWidget(theme, "Guard zone", 110, func(gtx layout.Context) layout.Dimensions {
			return this.guardZoneDropdown.Layout(gtx, theme)
		}),
		widgets.NewVerticalSpacerWidget(4),
		widgets.NewLabeledRowWidget(theme, "Guard preset", 110, func(gtx layout.Context) layout.Dimensions {
			return this.guardDropdown.Layout(gtx, theme)
		}),
		widgets.NewLabeledRowWidget(theme, "Guard value", 110, widgets.NewTextboxWidget(theme, &this.guardValueEdit, "guard value")),
		widgets.NewVerticalSpacerWidget(4),
		widgets.NewLabeledRowWidget(theme, "Weekly +", 110, func(gtx layout.Context) layout.Dimensions {
			return this.weeklyDropdown.Layout(gtx, theme)
		}),
		widgets.NewLabeledRowWidget(theme, "Increment", 110, widgets.NewTextboxWidget(theme, &this.weeklyEdit, "0.15")),
		widgets.NewVerticalSpacerWidget(6),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.advancedBool, "Advanced options"),
	)
	if this.advancedBool.Value {
		rows = append(rows,
			widgets.NewLabeledRowWidget(theme, "Match group", 110, widgets.NewTextboxWidget(theme, &this.matchGroupEdit, "rnd_guard_…")),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.escapeBool, "Guard escape"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.simSquadBool, "Sim turn squad"),
		)
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(10),
		widgets.NewButtonWidget(theme, "Delete this connection", &this.sidePropDelete, false),
	)
	return rows
}

// syncPropsFromConnection loads the property widgets from the selected connection.
// Called once whenever the selection changes.
func (this *ConnectionEditorDialog) syncPropsFromConnection() {
	connection := this.selected
	if connection == nil {
		return
	}
	this.typeDropdown.SetItems(connection_editor.UserCreatableConnectionTypes())
	if !this.typeDropdown.SelectByName(connection.ConnectionType) {
		this.typeDropdown.SelectByName("Direct")
	}
	this.guardZoneDropdown.SetItems([]string{connection.From, connection.To})
	if !this.guardZoneDropdown.SelectByName(connection.GuardZone) {
		this.guardZoneDropdown.SelectByName(connection.From)
	}
	tier := connection_editor.HigherTierOf(connection.From, connection.To, this.zones, this.playerZones)
	labels, values := guardPresetItems(tier)
	this.guardPresetValues = values
	this.guardDropdown.SetItems(labels)
	this.guardDropdown.SelectByName(matchGuardLabel(labels, values, connection.GuardValue))
	this.guardValueEdit.SetText(strconv.Itoa(connection.GuardValue))
	this.weeklyDropdown.SetItems(connection_editor.WeeklyIncrementLabels)
	this.weeklyDropdown.SelectByName(matchWeeklyLabel(connection.GuardWeeklyIncrement))
	this.weeklyEdit.SetText(formatIncrement(connection.GuardWeeklyIncrement))
	this.matchGroupEdit.SetText(connection.GuardMatchGroup)
	this.escapeBool.Value = connection.GuardEscape
	this.simSquadBool.Value = connection.SimTurnSquad
}

// writebackProps copies the property widget state back into the selected
// connection after the panel has been laid out for this frame.
func (this *ConnectionEditorDialog) writebackProps(gtx layout.Context) {
	connection := this.selected
	if connection == nil {
		return
	}
	if this.sidePropDelete.Clicked(gtx) {
		this.deleteConnection(connection)
		return
	}
	typeItems := connection_editor.UserCreatableConnectionTypes()
	if index := this.typeDropdown.GetSelectedIndex(); index >= 0 && index < len(typeItems) {
		connection.ConnectionType = typeItems[index]
	}
	zoneItems := []string{connection.From, connection.To}
	if index := this.guardZoneDropdown.GetSelectedIndex(); index >= 0 && index < len(zoneItems) {
		connection.GuardZone = zoneItems[index]
	}
	if this.guardDropdown.WasUpdated {
		if index := this.guardDropdown.GetSelectedIndex(); index >= 0 && index < len(this.guardPresetValues) {
			this.guardValueEdit.SetText(strconv.Itoa(this.guardPresetValues[index]))
		}
	}
	if value, err := strconv.Atoi(strings.TrimSpace(this.guardValueEdit.Text())); err == nil {
		connection.GuardValue = value
	}
	if this.weeklyDropdown.WasUpdated {
		if index := this.weeklyDropdown.GetSelectedIndex(); index >= 0 && index < len(connection_editor.WeeklyIncrementValues) {
			this.weeklyEdit.SetText(formatIncrement(connection_editor.WeeklyIncrementValues[index]))
		}
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.weeklyEdit.Text()), 64); err == nil {
		connection.GuardWeeklyIncrement = value
	}
	connection.GuardMatchGroup = strings.TrimSpace(this.matchGroupEdit.Text())
	connection.GuardEscape = this.escapeBool.Value
	connection.SimTurnSquad = this.simSquadBool.Value
}

func (this *ConnectionEditorDialog) addConnection(from, to string) {
	connection := connection_editor.NewDefaultConnection(from, to, this.zones, this.playerZones)
	this.working = append(this.working, &connection)
	this.selected = &connection
	this.syncedFor = nil
}

func (this *ConnectionEditorDialog) deleteConnection(connection *template.Connection) {
	for i, candidate := range this.working {
		if candidate == connection {
			this.working = append(this.working[:i], this.working[i+1:]...)
			break
		}
	}
	if this.selected == connection {
		this.selected = nil
		this.syncedFor = nil
	}
}

func (this *ConnectionEditorDialog) resetToOriginal() {
	this.working = this.working[:0]
	for i := range this.original {
		clone := connection_editor.CloneConnection(this.original[i], false)
		this.working = append(this.working, &clone)
	}
	this.selected = nil
	this.syncedFor = nil
	this.addMode = false
	this.pendingFrom = ""
	this.dragging = false
}

func derefConnections(pointers []*template.Connection) []template.Connection {
	out := make([]template.Connection, len(pointers))
	for i, pointer := range pointers {
		out[i] = *pointer
	}
	return out
}

func guardPresetItems(tier connection_editor.ZoneTier) (labels []string, values []int) {
	for _, extra := range connection_editor.ExtrasForTier(tier) {
		labels = append(labels, fmt.Sprintf("%s (%d)", extra.Label, extra.Value))
		values = append(values, extra.Value)
	}
	presets := connection_editor.GuardPresetsForTier(tier)
	for i, strength := range connection_editor.StrengthLabels {
		labels = append(labels, fmt.Sprintf("%s (%d)", strength, presets[i]))
		values = append(values, presets[i])
	}
	return labels, values
}

func matchGuardLabel(labels []string, values []int, value int) string {
	for i, candidate := range values {
		if candidate == value {
			return labels[i]
		}
	}
	return ""
}

func matchWeeklyLabel(value float64) string {
	for i, candidate := range connection_editor.WeeklyIncrementValues {
		if math.Abs(candidate-value) < 1e-9 {
			return connection_editor.WeeklyIncrementLabels[i]
		}
	}
	return ""
}

func formatIncrement(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

func drawCanvasText(gtx layout.Context, theme *material.Theme, center image.Point, text string, sizeSp int, textColor color.NRGBA) {
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
