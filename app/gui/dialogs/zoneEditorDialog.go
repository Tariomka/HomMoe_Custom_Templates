package dialogs

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

	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

var (
	colorEdgeSelected = themes.ColorEditorEdgeSelected
	colorUserAddedDot = themes.ColorEditorUserAddedDot
	colorGuardLabel   = themes.ColorEditorGuardLabel
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

// ZoneEditorDialog is the Manual Zone Editor. It renders zones as
// tier-coloured nodes and connections as curved lines, supports drag-to-create
// and right-click-delete for connections, lets the user add, move and delete
// zones, and exposes a property panel for the selected edge or zone. It works
// on private copies of the zone and connection lists; Apply commits, Cancel/✕
// discards.
type ZoneEditorDialog struct {
	zones         []template.Zone
	originalZones []template.Zone
	playerZones   map[string]bool
	topology      config.MapTopology
	tuning        models.GenerationTuning
	generateRoads bool
	working       []*template.Connection
	original      []template.Connection
	onApply       func([]template.Zone, []template.Connection)

	// Geometry recomputed every frame from BuildPreviewLayout.
	positions    map[string]image.Point
	previewZones []services.PreviewZone
	radius       int
	side         int
	edges        []connEdgeGeom

	// Canvas interaction.
	canvasTag     int
	selected      *template.Connection
	selectedZone  string
	addMode       bool
	pendingFrom   string
	dragging      bool
	dragPos       image.Point
	zoneDragName  string
	zoneDragMoved bool
	pressPos      image.Point
	hint          string

	// Toolbar / footer.
	addBtn     widget.Clickable
	addZoneBtn widget.Clickable
	deleteBtn  widget.Clickable
	resetBtn   widget.Clickable
	applyBtn   widget.Clickable
	cancelBtn  widget.Clickable

	// Property panel.
	sideScroll        widget.List
	syncedFor         *template.Connection
	typeDropdown      *components.DropdownSelector
	guardZoneDropdown *components.DropdownSelector
	guardDropdown     *components.DropdownSelector
	guardPresetValues []int
	weeklyDropdown    *components.DropdownSelector
	guardValueEdit    widget.Editor
	weeklyEdit        widget.Editor
	matchGroupEdit    widget.Editor
	advancedBool      widget.Bool
	escapeBool        widget.Bool
	simSquadBool      widget.Bool
	sidePropDelete    widget.Clickable

	// Zone property panel.
	syncedZoneFor   string
	qualityDropdown *components.DropdownSelector
	castleDropdown  *components.DropdownSelector
	zoneSizeEdit    widget.Editor
	zoneGuardEdit   widget.Editor
	zoneWeeklyEdit  widget.Editor
	sideZoneDelete  widget.Clickable
}

// NewZoneEditorDialog builds an editor over copies of the given zones and
// connections. onApply receives the edited zone and connection lists when the
// user applies.
func NewZoneEditorDialog(
	zones []template.Zone,
	connections []template.Connection,
	topology config.MapTopology,
	tuning models.GenerationTuning,
	generateRoads bool,
	onApply func([]template.Zone, []template.Connection)) *ZoneEditorDialog {
	players := make(map[string]bool)
	for _, zone := range zones {
		if strings.HasPrefix(zone.Name, "Spawn-") {
			players[zone.Name] = true
		}
	}

	dialog := &ZoneEditorDialog{
		zones:             append([]template.Zone(nil), zones...),
		originalZones:     append([]template.Zone(nil), zones...),
		playerZones:       players,
		topology:          topology,
		tuning:            tuning,
		generateRoads:     generateRoads,
		onApply:           onApply,
		typeDropdown:      components.NewDropdownSelector(connection_editor.UserCreatableConnectionTypes()),
		guardZoneDropdown: components.NewDropdownSelector(nil),
		guardDropdown:     components.NewDropdownSelector(nil),
		weeklyDropdown:    components.NewDropdownSelector(connection_editor.WeeklyIncrementLabels),
		qualityDropdown:   components.NewDropdownSelector(connection_editor.QualityLabels),
		castleDropdown:    components.NewDropdownSelector([]string{"0", "1", "2", "3", "4"}),
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
	dialog.zoneSizeEdit.SingleLine = true
	dialog.zoneGuardEdit.SingleLine = true
	dialog.zoneWeeklyEdit.SingleLine = true
	return dialog
}

func (this *ZoneEditorDialog) Title() string { return "Manual Zone Editor" }

func (this *ZoneEditorDialog) PreferredSize() (unit.Dp, unit.Dp) {
	return unit.Dp(1000), unit.Dp(720)
}

func (this *ZoneEditorDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.applyBtn.Clicked(gtx) {
		if this.onApply != nil {
			this.onApply(this.zones, derefConnections(this.working))
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
		this.hint = ""
	}
	if this.addZoneBtn.Clicked(gtx) {
		this.addZone()
	}
	if this.deleteBtn.Clicked(gtx) {
		if this.selected != nil {
			this.deleteConnection(this.selected)
		} else if this.selectedZone != "" {
			this.deleteZone(this.selectedZone)
		}
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

func (this *ZoneEditorDialog) layoutToolbar(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		addLabel := "Add connection"
		if this.addMode {
			addLabel = "Adding… (click empty to cancel)"
		}
		hasSelection := this.selected != nil || this.selectedZone != ""
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(widgets.NewToggleButtonWidget(theme, addLabel, &this.addBtn, this.addMode)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewToggleButtonWidget(theme, "Add zone", &this.addZoneBtn, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Delete selected", &this.deleteBtn, !hasSelection)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Reset to generated", &this.resetBtn, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(12)),
			layout.Flexed(1, this.layoutStatus(theme)),
		)
	}
}

func (this *ZoneEditorDialog) layoutStatus(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		connections := derefConnections(this.working)
		if connection_editor.ComputeHasErrors(this.zones, connections) {
			label := material.Body2(theme, "⚠ A connection references a missing zone — fix before export.")
			label.Color = themes.ColorError
			label.TextSize = unit.Sp(12)
			label.MaxLines = 2
			return label.Layout(gtx)
		}
		message := fmt.Sprintf("%d zones · %d connections", len(this.zones), len(connections))
		if this.hint != "" {
			message = this.hint
		} else if this.addMode {
			message = "Add mode: press a zone and drag to another to connect."
		} else if isolated := connection_editor.FindIsolatedZones(this.zones, connections); len(isolated) > 0 {
			message = fmt.Sprintf("%d zones · %d connections · %d isolated zone(s)", len(this.zones), len(connections), len(isolated))
		}
		label := material.Body2(theme, message)
		label.Color = themes.ColorTextDim
		label.TextSize = unit.Sp(12)
		label.MaxLines = 2
		return label.Layout(gtx)
	}
}

func (this *ZoneEditorDialog) layoutFooter(theme *material.Theme) layout.Widget {
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
func (this *ZoneEditorDialog) layoutCanvas(gtx layout.Context, theme *material.Theme) layout.Dimensions {
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
		this.addConnection(from, target)
		this.addMode = false
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

func (this *ZoneEditorDialog) hitTestEdge(pos image.Point) *template.Connection {
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
func (this *ZoneEditorDialog) recomputeGeometry(side int) {
	this.side = side
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
	paint.FillShape(gtx.Ops, colorEdgeSelected, clip.Stroke{Path: path.End(), Width: float32(gtx.Dp(unit.Dp(2)))}.Op())
}

func (this *ZoneEditorDialog) drawNodes(gtx layout.Context, theme *material.Theme) {
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
	if !this.addMode && this.selectedZone != "" {
		if center, ok := this.positions[this.selectedZone]; ok {
			reach := this.radius + 4
			rect := image.Rect(center.X-reach, center.Y-reach, center.X+reach, center.Y+reach)
			paint.FillShape(gtx.Ops, colorEdgeSelected, clip.Stroke{
				Path:  clip.UniformRRect(rect, reach).Path(gtx.Ops),
				Width: float32(gtx.Dp(unit.Dp(2))),
			}.Op())
		}
	}
}

func (this *ZoneEditorDialog) layoutSidePanel(gtx layout.Context, theme *material.Theme) layout.Dimensions {
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
		if zone := this.selectedZoneRef(); zone != nil {
			if this.syncedZoneFor != zone.Name {
				this.syncZoneProps(zone)
				this.syncedZoneFor = zone.Name
			}
			rows := this.zonePropertyRows(theme, zone)
			dims := material.List(theme, &this.sideScroll).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
				return rows[index](gtx)
			})
			this.writebackZoneProps(gtx, zone)
			return dims
		}
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

func (this *ZoneEditorDialog) layoutSideHint(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, "Nothing selected")
			label.Color = themes.ColorText
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Click a zone to edit its size, quality and guards — drag it to move it.")),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Click a connection line to edit its guard, type and growth.")),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Use “Add connection”, then drag from one zone to another to create a link.")),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(widgets.NewDimmedLabelWidget(theme, "Use “Add zone” to add a neutral zone. Right-click a line to delete it.")),
	)
}

func (this *ZoneEditorDialog) propertyRows(theme *material.Theme) []layout.Widget {
	connection := this.selected
	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, connection.From+"  →  "+connection.To)
			label.Color = themes.ColorAccentBright
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
func (this *ZoneEditorDialog) syncPropsFromConnection() {
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
func (this *ZoneEditorDialog) writebackProps(gtx layout.Context) {
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

func (this *ZoneEditorDialog) addConnection(from, to string) {
	connection := connection_editor.NewDefaultConnection(from, to, this.zones, this.playerZones)
	this.working = append(this.working, &connection)
	this.selected = &connection
	this.syncedFor = nil
}

func (this *ZoneEditorDialog) deleteConnection(connection *template.Connection) {
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

func (this *ZoneEditorDialog) resetToOriginal() {
	this.zones = append([]template.Zone(nil), this.originalZones...)
	this.working = this.working[:0]
	for i := range this.original {
		clone := connection_editor.CloneConnection(this.original[i], false)
		this.working = append(this.working, &clone)
	}
	this.selected = nil
	this.syncedFor = nil
	this.selectedZone = ""
	this.syncedZoneFor = ""
	this.addMode = false
	this.pendingFrom = ""
	this.dragging = false
	this.zoneDragName = ""
	this.hint = ""
}

// selectZone makes the named zone the active selection in the side panel.
func (this *ZoneEditorDialog) selectZone(name string) {
	this.selectedZone = name
	this.selected = nil
	this.syncedFor = nil
	this.hint = ""
}

func (this *ZoneEditorDialog) selectedZoneRef() *template.Zone {
	if this.selectedZone == "" {
		return nil
	}
	return this.zoneByName(this.selectedZone)
}

func (this *ZoneEditorDialog) zoneByName(name string) *template.Zone {
	for i := range this.zones {
		if this.zones[i].Name == name {
			return &this.zones[i]
		}
	}
	return nil
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
	x := math.Min(math.Max(float64(pos.X)/float64(this.side), 0.04), 0.96)
	y := math.Min(math.Max(float64(pos.Y)/float64(this.side), 0.04), 0.96)
	zone.ManualPosition = &[2]float64{x, y}
}

// ensureManualPositions freezes the current generated layout into manual
// positions on every zone, so the preview keeps every other zone in place
// while one is moved or added.
func (this *ZoneEditorDialog) ensureManualPositions() {
	if this.side <= 0 {
		return
	}
	for i := range this.zones {
		if this.zones[i].ManualPosition != nil {
			continue
		}
		if pos, ok := this.positions[this.zones[i].Name]; ok {
			this.zones[i].ManualPosition = &[2]float64{
				float64(pos.X) / float64(this.side),
				float64(pos.Y) / float64(this.side),
			}
		} else {
			open := connection_editor.FindOpenPosition(this.manualPositions())
			this.zones[i].ManualPosition = &open
		}
	}
}

func (this *ZoneEditorDialog) manualPositions() [][2]float64 {
	out := make([][2]float64, 0, len(this.zones))
	for _, zone := range this.zones {
		if zone.ManualPosition != nil {
			out = append(out, *zone.ManualPosition)
		}
	}
	return out
}

// addZone appends a new medium-quality neutral zone at an open spot.
func (this *ZoneEditorDialog) addZone() {
	label := connection_editor.NextFreeZoneLabel(this.zones)
	if label == "" {
		this.hint = "Zone label pool exhausted — cannot add more zones."
		return
	}
	this.ensureManualPositions()
	zone := connection_editor.NewDefaultNeutralZone(label, models.QualityMedium, 1, this.generateRoads, this.tuning)
	position := connection_editor.FindOpenPosition(this.manualPositions())
	zone.ManualPosition = &position
	this.zones = append(this.zones, zone)
	this.addMode = false
	this.pendingFrom = ""
	this.selectZone(zone.Name)
	this.syncedZoneFor = ""
	this.hint = fmt.Sprintf("Added %s — drag it into place, then connect it with “Add connection”.", zone.Name)
}

// deleteZone removes a zone and every connection referencing it. Spawn zones
// are protected; player count is managed in the General tab.
func (this *ZoneEditorDialog) deleteZone(name string) {
	if !connection_editor.CanDeleteZone(name, this.playerZones) {
		this.hint = "Spawn zones can't be deleted — change player count in the General tab."
		return
	}
	zones, connections := connection_editor.RemoveZone(this.zones, derefConnections(this.working), name)
	this.zones = zones
	this.working = this.working[:0]
	for i := range connections {
		kept := connections[i]
		this.working = append(this.working, &kept)
	}
	this.selected = nil
	this.syncedFor = nil
	if this.selectedZone == name {
		this.selectedZone = ""
		this.syncedZoneFor = ""
	}
	this.hint = ""
}

func (this *ZoneEditorDialog) zonePropertyRows(theme *material.Theme, zone *template.Zone) []layout.Widget {
	isNeutral := strings.HasPrefix(zone.Name, "Neutral-")
	isSpawn := this.playerZones[zone.Name]
	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, zone.Name)
			label.Color = themes.ColorAccentBright
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		},
	}
	if isSpawn {
		rows = append(rows, widgets.NewDimmedLabelWidget(theme, "Player spawn zone — content is managed by the generator."))
	} else if !isNeutral {
		rows = append(rows, widgets.NewDimmedLabelWidget(theme, "Quality presets apply to neutral zones only."))
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(6),
		widgets.NewLabeledRowWidget(theme, "Size", 110, widgets.NewTextboxWidget(theme, &this.zoneSizeEdit, "0.1 – 2.0")),
		widgets.NewVerticalSpacerWidget(4),
		widgets.NewLabeledRowWidget(theme, "Guard ×", 110, widgets.NewTextboxWidget(theme, &this.zoneGuardEdit, "guard multiplier")),
		widgets.NewLabeledRowWidget(theme, "Weekly +", 110, widgets.NewTextboxWidget(theme, &this.zoneWeeklyEdit, "0.15")),
	)
	if isNeutral {
		rows = append(rows,
			widgets.NewVerticalSpacerWidget(4),
			widgets.NewLabeledRowWidget(theme, "Quality", 110, func(gtx layout.Context) layout.Dimensions {
				return this.qualityDropdown.Layout(gtx, theme)
			}),
			widgets.NewLabeledRowWidget(theme, "Castles", 110, func(gtx layout.Context) layout.Dimensions {
				return this.castleDropdown.Layout(gtx, theme)
			}),
			widgets.NewDimmedLabelWidget(theme, "Changing quality or castles regenerates the zone's content."),
		)
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(10),
		widgets.NewButtonWidget(theme, "Delete this zone", &this.sideZoneDelete, isSpawn),
	)
	return rows
}

// syncZoneProps loads the zone property widgets from the selected zone.
// Called once whenever the zone selection changes.
func (this *ZoneEditorDialog) syncZoneProps(zone *template.Zone) {
	quality := connection_editor.QualityOfZone(*zone)
	this.qualityDropdown.SelectByName(connection_editor.QualityLabels[int(quality)])
	castles := connection_editor.CountZoneCastles(*zone)
	if castles > 4 {
		castles = 4
	}
	this.castleDropdown.SelectByName(strconv.Itoa(castles))
	this.zoneSizeEdit.SetText(strconv.FormatFloat(zone.Size, 'f', -1, 64))
	this.zoneGuardEdit.SetText(strconv.FormatFloat(zone.GuardMultiplier, 'f', -1, 64))
	this.zoneWeeklyEdit.SetText(formatIncrement(zone.GuardWeeklyIncrement))
}

// writebackZoneProps copies the zone widget state back into the selected zone
// after the panel has been laid out for this frame.
func (this *ZoneEditorDialog) writebackZoneProps(gtx layout.Context, zone *template.Zone) {
	if this.sideZoneDelete.Clicked(gtx) {
		this.deleteZone(zone.Name)
		return
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.zoneSizeEdit.Text()), 64); err == nil {
		zone.Size = math.Round(math.Min(math.Max(value, 0.1), 2.0)*100) / 100
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.zoneGuardEdit.Text()), 64); err == nil {
		zone.GuardMultiplier = value
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.zoneWeeklyEdit.Text()), 64); err == nil {
		zone.GuardWeeklyIncrement = value
	}
	if strings.HasPrefix(zone.Name, "Neutral-") && (this.qualityDropdown.WasUpdated || this.castleDropdown.WasUpdated) {
		quality := models.NeutralZoneQuality(this.qualityDropdown.GetSelectedIndex())
		castles := this.castleDropdown.GetSelectedIndex()
		connection_editor.ApplyNeutralZoneQuality(zone, quality, castles, this.tuning)
		this.syncedZoneFor = "" // re-sync dependent fields next frame
	}
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
