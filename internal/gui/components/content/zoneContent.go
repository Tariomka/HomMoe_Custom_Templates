package content

import (
	"fmt"
	"iter"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// zoneContentRow is one editable item inside a zone-content section.
type zoneContentRow struct {
	Mapping     models.SidMapping
	Count       int
	IsGuarded   widget.Bool
	NearCastle  widget.Bool
	RoadDistIdx int
	IsGroup     bool

	countSld  widget.Float
	roadCombo *DropdownSelector
	removeBtn widget.Clickable
	dupBtn    widget.Clickable
}

func newZoneContentRow(mapping models.SidMapping, count int, guarded, nearCastle bool, roadIdx int, isGroup bool) *zoneContentRow {
	row := &zoneContentRow{
		Mapping:     mapping,
		Count:       count,
		RoadDistIdx: roadIdx,
		IsGroup:     isGroup,
		roadCombo:   NewDropdownSelector(constants.RoadDistances),
	}
	if row.RoadDistIdx >= 0 && row.RoadDistIdx < len(constants.RoadDistances) {
		row.roadCombo.selectedIndex = row.RoadDistIdx
	}
	row.IsGuarded.Value = guarded
	row.NearCastle.Value = nearCastle
	return row
}

// ZoneContentSection is one of the four mandatory-content groups.
type ZoneContentSection struct {
	Title     string
	Items     []models.SidMapping
	MaxCount  int
	ShowNear  bool
	rows      []*zoneContentRow
	addPreset *DropdownSelector
	addBtn    widget.Clickable
}

func NewZoneContentSection(title string, items []models.SidMapping, maxCount int, showNear bool) *ZoneContentSection {
	labels := make([]string, len(items))
	for i, item := range items {
		labels[i] = item.Name
	}
	return &ZoneContentSection{
		Title:     title,
		Items:     items,
		MaxCount:  maxCount,
		ShowNear:  showNear,
		addPreset: NewDropdownSelector(labels),
	}
}

// Add appends a new row using the given mapping with sensible defaults.
func (this *ZoneContentSection) Add(mapping models.SidMapping, count int, guarded, near bool, roadIdx int, group bool) {
	if count < 1 {
		count = 1
	}
	if count > this.MaxCount {
		count = this.MaxCount
	}
	this.rows = append(this.rows, newZoneContentRow(mapping, count, guarded, near, roadIdx, group))
}

func (this *ZoneContentSection) ClearRows() {
	this.rows = nil
}

func (this *ZoneContentSection) IterateRows() iter.Seq[*zoneContentRow] {
	return func(yield func(*zoneContentRow) bool) {
		for _, row := range this.rows {
			if !yield(row) {
				return
			}
		}
	}
}

func (this *ZoneContentSection) Layout(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Process per-section button clicks first.
		if this.addBtn.Clicked(gtx) && len(this.Items) > 0 {
			idx := this.addPreset.GetSelectedIndex()
			if idx < 0 || idx >= len(this.Items) {
				idx = 0
			}
			this.Add(this.Items[idx], 1, true, false, 0, false)
		}
		// Process per-row clicks (collect indices to remove).
		keep := this.rows[:0]
		for i, row := range this.rows {
			if row.removeBtn.Clicked(gtx) {
				continue
			}
			if row.dupBtn.Clicked(gtx) {
				keep = append(keep, row)
				clone := newZoneContentRow(row.Mapping, row.Count, row.IsGuarded.Value, row.NearCastle.Value, row.RoadDistIdx, row.IsGroup)
				keep = append(keep, clone)
				continue
			}
			_ = i
			keep = append(keep, row)
		}
		this.rows = keep

		return widgets.NewSectionWidget(theme, this.Title, []layout.Widget{
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(120)
						label := material.Body2(theme, "Add preset:")
						label.Color = themes.ColorTextDim
						label.TextSize = unit.Sp(12)
						return label.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return this.addPreset.Layout(gtx, theme)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(widgets.NewButtonWidget(theme, "+ Add", &this.addBtn, false)),
				)
			},
			func(gtx layout.Context) layout.Dimensions {
				if len(this.rows) == 0 {
					label := material.Body2(theme, "(no items)")
					label.Color = themes.ColorTextDim
					label.TextSize = unit.Sp(12)
					return layout.Inset{Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, label.Layout)
				}
				children := make([]layout.FlexChild, 0, len(this.rows)*2)
				for i, row := range this.rows {
					if i > 0 {
						children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
					}
					children = append(children, layout.Rigid(this.layoutRow(theme, row)))
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			},
		})(gtx)
	}
}

func (this *ZoneContentSection) layoutRow(theme *material.Theme, row *zoneContentRow) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Sync count slider into integer field.
		desired := utils.Normalize(float32(row.Count), 1, float32(this.MaxCount))
		if !row.countSld.Dragging() && row.countSld.Value == 0 && row.Count > 0 {
			row.countSld.Value = desired
		}
		liveCount := utils.RoundedRange(row.countSld.Value, 1, this.MaxCount)
		row.Count = liveCount
		row.RoadDistIdx = row.roadCombo.GetSelectedIndex()

		return widgets.NewPanelWidget(unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(theme, row.Mapping.Name)
							label.Color = themes.ColorGold
							label.TextSize = unit.Sp(13)
							return label.Layout(gtx)
						}),
						layout.Rigid(widgets.NewButtonWidget(theme, "Duplicate", &row.dupBtn, false)),
						layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
						layout.Rigid(widgets.NewButtonWidget(theme, "Remove", &row.removeBtn, false)),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(widgets.NewLabeledRowWidget(theme, "Count", 100, widgets.NewLabeledSlider(theme, &row.countSld, fmt.Sprintf("%d", liveCount)))),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &row.IsGuarded, "Guarded")),
						layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !this.ShowNear {
								return layout.Dimensions{}
							}
							return widgets.NewLabeledCheckboxRowWidget(theme, &row.NearCastle, "Near castle")(gtx)
						}),
					)
				}),
				layout.Rigid(widgets.NewLabeledRowWidget(theme, "Road distance", 100, func(gtx layout.Context) layout.Dimensions {
					return row.roadCombo.Layout(gtx, theme)
				})),
			)
		})(gtx)
	}
}
