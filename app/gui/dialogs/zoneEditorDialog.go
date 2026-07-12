package dialogs

import (
	"fmt"
	"image"
	"math"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

// ZoneEditorDialog is the Manual Zone Editor. It renders zones as
// tier-coloured nodes and connections as curved lines, supports drag-to-create
// and right-click-delete for connections, lets the user add, move and delete
// zones, and exposes a property panel for the selected edge or zone. It works
// on private copies of the zone and connection lists; Apply commits, Cancel/✕
// discards.
type ZoneEditorDialog struct {
	zones         []entities.Zone
	originalZones []entities.Zone
	playerZones   map[string]bool
	topology      config.MapTopology
	tuning        models.GenerationTuning
	generateRoads bool
	working       []*entities.Connection
	original      []entities.Connection
	onApply       func([]entities.Zone, []entities.Connection)
	layoutService *preview_service.PreviewLayoutService

	// Geometry recomputed every frame from BuildPreviewLayout.
	positions    map[string]image.Point
	previewZones []preview.Zone
	radius       int
	side         int
	edges        []connEdgeGeom

	// Canvas interaction.
	canvasTag     int
	selected      *entities.Connection
	selectedZone  string
	addMode       bool
	addZoneMode   bool
	pendingFrom   string
	dragging      bool
	dragPos       image.Point
	zoneDragName  string
	zoneDragMoved bool
	pressPos      image.Point
	hint          string

	// Active zone-alignment snap guides (canvas px), valid while a zone is
	// being dragged and holding onto another zone's edge/centre extension.
	snapGuideX       float64
	snapGuideY       float64
	snapGuideXActive bool
	snapGuideYActive bool

	// Toolbar / footer.
	addBtn     widget.Clickable
	addZoneBtn widget.Clickable
	deleteBtn  widget.Clickable
	resetBtn   widget.Clickable
	snapBool   widget.Bool
	applyBtn   widget.Clickable
	cancelBtn  widget.Clickable

	// Property panel.
	sideScroll        widget.List
	syncedFor         *entities.Connection
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
	zones []entities.Zone,
	connections []entities.Connection,
	topology config.MapTopology,
	tuning models.GenerationTuning,
	generateRoads bool,
	onApply func([]entities.Zone, []entities.Connection)) *ZoneEditorDialog {
	players := make(map[string]bool)
	for _, zone := range zones {
		if strings.HasPrefix(zone.Name, "Spawn-") {
			players[zone.Name] = true
		}
	}

	dialog := &ZoneEditorDialog{
		zones:             append([]entities.Zone(nil), zones...),
		originalZones:     append([]entities.Zone(nil), zones...),
		playerZones:       players,
		topology:          topology,
		tuning:            tuning,
		generateRoads:     generateRoads,
		onApply:           onApply,
		layoutService:     preview_service.NewPreviewLayoutService(),
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
		this.addZoneMode = false
		this.pendingFrom = ""
		this.dragging = false
		this.hint = ""
	}
	if this.addZoneBtn.Clicked(gtx) {
		this.addZoneMode = !this.addZoneMode
		this.addMode = false
		this.pendingFrom = ""
		this.dragging = false
		this.hint = ""
	}
	if this.deleteBtn.Clicked(gtx) {
		if this.selected != nil {
			this.deleteConnection(this.selected)
		} else if this.selectedZone != "" {
			this.deleteZone(this.selectedZone)
		}
	}
	// The side-panel delete buttons must be polled BEFORE their Clickables are
	// laid out - Clickable.Layout consumes the click, so a check after layout
	// never fires.
	if this.sidePropDelete.Clicked(gtx) && this.selected != nil {
		this.deleteConnection(this.selected)
	}
	if this.sideZoneDelete.Clicked(gtx) {
		if zone := this.selectedZoneRef(); zone != nil {
			this.deleteZone(zone.Name)
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
			addLabel = "Adding... (click empty to stop)"
		}
		addZoneLabel := "Add zone"
		if this.addZoneMode {
			addZoneLabel = "Placing... (click a zone to stop)"
		}
		hasSelection := this.selected != nil || this.selectedZone != ""
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(widgets.NewToggleButtonWidget(theme, addLabel, &this.addBtn, this.addMode)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewToggleButtonWidget(theme, addZoneLabel, &this.addZoneBtn, this.addZoneMode)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Delete selected", &this.deleteBtn, !hasSelection)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Reset to generated", &this.resetBtn, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(10)),
			layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.snapBool, "Snap")),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(12)),
			layout.Flexed(1, this.layoutStatus(theme)),
		)
	}
}

func (this *ZoneEditorDialog) layoutStatus(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		connections := derefConnections(this.working)
		if connection_editor.ComputeHasErrors(this.zones, connections) {
			label := material.Body2(theme, "⚠ A connection references a missing zone - fix before export.")
			label.Color = themes.ColorError
			label.TextSize = unit.Sp(12)
			label.MaxLines = 2
			return label.Layout(gtx)
		}
		message := fmt.Sprintf("%d zones · %d connections", len(this.zones), len(connections))
		if this.hint != "" {
			message = this.hint
		} else if this.addMode {
			message = "Add mode: press a zone and drag to another to connect. Repeat to add more."
		} else if this.addZoneMode {
			message = "Add zone mode: click an empty spot to place a zone. Repeat to add more."
		} else if isolated := connection_editor.FindIsolatedZones(this.zones, connections); len(isolated) > 0 {
			message = fmt.Sprintf(
				"%d zones · %d connections · %d isolated zone(s)",
				len(this.zones),
				len(connections),
				len(isolated),
			)
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
			layout.Rigid(widgets.NewBrightButtonWidget(theme, "Apply changes", &this.applyBtn, false)),
		)
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
			dims := material.List(theme, &this.sideScroll).
				Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
					return rows[index](gtx)
				})
			this.writebackZoneProps(zone)
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
		dims := material.List(theme, &this.sideScroll).
			Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
				return rows[index](gtx)
			})
		this.writebackProps()
		return dims
	})
	return layout.Dimensions{Size: size}
}

func (this *ZoneEditorDialog) layoutSideHint(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, "Nothing selected")
			label.Color = themes.ColorText
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Rigid(
			widgets.NewDimmedLabelWidget(
				theme,
				"Click a zone to edit its size, quality and guards - drag it to move it.",
			),
		),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(
			widgets.NewDimmedLabelWidget(theme, "Click a connection line to edit its guard, type and growth."),
		),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(
			widgets.NewDimmedLabelWidget(
				theme,
				"Use “Add connection”, then drag from one zone to another to create a link.",
			),
		),
		layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
		layout.Rigid(
			widgets.NewDimmedLabelWidget(
				theme,
				"Use “Add zone”, then click empty spots to place neutral zones. Right-click a line to delete it.",
			),
		),
	)
}

func (this *ZoneEditorDialog) addConnection(from, to string) {
	connection := connection_editor.NewDefaultConnection(from, to, this.zones, this.playerZones)
	this.working = append(this.working, &connection)
	this.selected = &connection
	this.syncedFor = nil
}

func (this *ZoneEditorDialog) deleteConnection(connection *entities.Connection) {
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
	this.zones = append([]entities.Zone(nil), this.originalZones...)
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
	this.addZoneMode = false
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

func (this *ZoneEditorDialog) selectedZoneRef() *entities.Zone {
	if this.selectedZone == "" {
		return nil
	}
	return this.zoneByName(this.selectedZone)
}

func (this *ZoneEditorDialog) zoneByName(name string) *entities.Zone {
	for i := range this.zones {
		if this.zones[i].Name == name {
			return &this.zones[i]
		}
	}
	return nil
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

// addZoneAt appends a new medium-quality neutral zone at the clicked canvas
// position. The add-zone mode stays active so several zones can be placed
// without re-clicking the toolbar button.
func (this *ZoneEditorDialog) addZoneAt(pos image.Point) {
	if this.side <= 0 {
		return
	}
	label := connection_editor.NextFreeZoneLabel(this.zones)
	if label == "" {
		this.hint = "Zone label pool exhausted - cannot add more zones."
		return
	}
	this.ensureManualPositions()
	zone := connection_editor.NewDefaultNeutralZone(
		label, neutralZone.QualityMedium, 1, this.generateRoads, this.tuning)
	x := math.Min(math.Max(float64(pos.X)/float64(this.side), 0.04), 0.96)
	y := math.Min(math.Max(float64(pos.Y)/float64(this.side), 0.04), 0.96)
	zone.ManualPosition = &[2]float64{x, y}
	this.zones = append(this.zones, zone)
	this.selectZone(zone.Name)
	this.syncedZoneFor = ""
	this.hint = fmt.Sprintf("Added %s - connect it with “Add connection”.", zone.Name)
}

// deleteZone removes a zone and every connection referencing it. Spawn zones
// are protected; player count is managed in the General tab.
func (this *ZoneEditorDialog) deleteZone(name string) {
	if !connection_editor.CanDeleteZone(name, this.playerZones) {
		this.hint = "Spawn zones can't be deleted - change player count in the General tab."
		return
	}
	zones, connections := connection_editor.RemoveZone(this.zones, derefConnections(this.working), name)
	this.zones = zones
	this.working = this.working[:0]
	for i := range connections {
		this.working = append(this.working, &connections[i])
	}
	this.selected = nil
	this.syncedFor = nil
	if this.selectedZone == name {
		this.selectedZone = ""
		this.syncedZoneFor = ""
	}
	this.hint = ""
}

func derefConnections(pointers []*entities.Connection) []entities.Connection {
	out := make([]entities.Connection, len(pointers))
	for i, pointer := range pointers {
		out[i] = *pointer
	}
	return out
}
