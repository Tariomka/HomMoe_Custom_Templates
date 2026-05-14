package gui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
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
	roadCombo *comboBox
	removeBtn widget.Clickable
	dupBtn    widget.Clickable
}

func newZoneContentRow(mapping models.SidMapping, count int, guarded, nearCastle bool, roadIdx int, isGroup bool) *zoneContentRow {
	row := &zoneContentRow{
		Mapping:     mapping,
		Count:       count,
		RoadDistIdx: roadIdx,
		IsGroup:     isGroup,
		roadCombo:   newComboBox(roadDistances),
	}
	if row.RoadDistIdx >= 0 && row.RoadDistIdx < len(roadDistances) {
		row.roadCombo.Selected = row.RoadDistIdx
	}
	row.IsGuarded.Value = guarded
	row.NearCastle.Value = nearCastle
	return row
}

// zoneContentSection is one of the four mandatory-content groups.
type zoneContentSection struct {
	Title     string
	Items     []models.SidMapping
	MaxCount  int
	ShowNear  bool
	rows      []*zoneContentRow
	addPreset *comboBox
	addBtn    widget.Clickable
}

func newZoneContentSection(title string, items []models.SidMapping, maxCount int, showNear bool) *zoneContentSection {
	labels := make([]string, len(items))
	for i, item := range items {
		labels[i] = item.Name
	}
	return &zoneContentSection{
		Title:     title,
		Items:     items,
		MaxCount:  maxCount,
		ShowNear:  showNear,
		addPreset: newComboBox(labels),
	}
}

// Add appends a new row using the given mapping with sensible defaults.
func (this *zoneContentSection) Add(mapping models.SidMapping, count int, guarded, near bool, roadIdx int, group bool) {
	if count < 1 {
		count = 1
	}
	if count > this.MaxCount {
		count = this.MaxCount
	}
	this.rows = append(this.rows, newZoneContentRow(mapping, count, guarded, near, roadIdx, group))
}

func (this *zoneContentSection) Layout(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Process per-section button clicks first.
		if this.addBtn.Clicked(gtx) && len(this.Items) > 0 {
			idx := this.addPreset.Selected
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

		return NewSectionWidget(theme, this.Title, []layout.Widget{
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(120)
						label := material.Body2(theme, "Add preset:")
						label.Color = colTextDim
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
					label.Color = colTextDim
					label.TextSize = unit.Sp(12)
					return layout.Inset{Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, label.Layout)
				}
				children := make([]layout.FlexChild, 0, len(this.rows)*2)
				for i, row := range this.rows {
					row := row
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

func (this *zoneContentSection) layoutRow(theme *material.Theme, row *zoneContentRow) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Sync count slider into integer field.
		desired := mapRangeInv(float32(row.Count), 1, float32(this.MaxCount))
		if !row.countSld.Dragging() && row.countSld.Value == 0 && row.Count > 0 {
			row.countSld.Value = desired
		}
		liveCount := roundedRange(row.countSld.Value, 1, this.MaxCount)
		row.Count = liveCount
		row.RoadDistIdx = row.roadCombo.Selected

		return widgets.NewPanelWidget(unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(theme, row.Mapping.Name)
							label.Color = colGold
							label.TextSize = unit.Sp(13)
							return label.Layout(gtx)
						}),
						layout.Rigid(widgets.NewButtonWidget(theme, "Duplicate", &row.dupBtn, false)),
						layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
						layout.Rigid(widgets.NewButtonWidget(theme, "Remove", &row.removeBtn, false)),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(widgets.NewLabeledRowWidget(theme, "Count", 100, func(gtx layout.Context) layout.Dimensions {
					return sliderLabeled(gtx, theme, &row.countSld, fmt.Sprintf("%d", liveCount))
				})),
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

// — State integration —

// seedDefaultPlayerZoneContent mirrors C# InitializeDefaultPlayerZoneContents.
func (this *State) seedDefaultPlayerZoneContent() {
	this.zcMines.rows = nil
	this.zcTreasures.rows = nil
	this.zcHires.rows = nil
	this.zcBanks.rows = nil

	// Mines: wood/ore/gold guarded next-to-castle; crystals/mercury/gemstones/alchemy-lab guarded near road.
	this.zcMines.Add(constants.ContentIds.MineWood, 1, true, true, 0, false)
	this.zcMines.Add(constants.ContentIds.MineOre, 1, true, true, 0, false)
	this.zcMines.Add(constants.ContentIds.MineGold, 1, true, true, 0, false)
	this.zcMines.Add(constants.ContentIds.MineCrystals, 1, true, false, 1, false)
	this.zcMines.Add(constants.ContentIds.MineMercury, 1, true, false, 1, false)
	this.zcMines.Add(constants.ContentIds.MineGemstones, 1, true, false, 1, false)
	this.zcMines.Add(constants.ContentIds.AlchemyLab, 1, true, false, 1, false)

	// Treasures: PandoraBox + RandomItemEpic guarded.
	this.zcTreasures.Add(constants.ContentIds.PandoraBox, 1, true, false, 0, false)
	this.zcTreasures.Add(constants.ContentIds.RandomItemEpic, 1, true, false, 0, false)

	// Random hires: low ×2, high ×1, all-tier ×1 (groups).
	this.zcHires.Add(constants.IncludeListIds.RandomHiresLowTier, 2, true, false, 0, true)
	this.zcHires.Add(constants.IncludeListIds.RandomHiresHighTier, 1, true, false, 0, true)
	this.zcHires.Add(constants.IncludeListIds.RandomHiresAllTier, 1, true, false, 0, true)

	// Resource banks: tier1 ×2, tier2 ×1.
	this.zcBanks.Add(constants.IncludeListIds.ResourceBanksTier1, 2, true, false, 0, true)
	this.zcBanks.Add(constants.IncludeListIds.ResourceBanksTier2, 1, true, false, 0, true)
}

// applyZoneContentItems replaces every section based on a flat list of items
// loaded from a settings file. Items are routed to the appropriate section by
// SID lookup.
func (this *State) applyZoneContentItems(items []models.ZoneContentItem) {
	this.zcMines.rows = nil
	this.zcTreasures.rows = nil
	this.zcHires.rows = nil
	this.zcBanks.rows = nil
	for _, item := range items {
		mapping := models.SidMapping{Sid: item.Sid, Name: item.Name}
		if found, ok := helpers.LookupSid(item.Sid); ok {
			mapping = found
		}
		count := max(item.Count, 1)
		roadIdx := max(indexOf(roadDistances, item.RoadDistance), 0)
		switch {
		case sectionContains(constants.ContentItemGroup.Mines, item.Sid):
			this.zcMines.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		case sectionContains(constants.ContentItemGroup.Treasures, item.Sid):
			this.zcTreasures.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		case sectionContains(constants.ContentItemGroup.HireBuildings, item.Sid):
			this.zcHires.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		case sectionContains(constants.ContentItemGroup.ResourceBanks, item.Sid):
			this.zcBanks.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		default:
			// Unknown SID — keep with treasures by default.
			this.zcTreasures.Add(mapping, count, item.IsGuarded, item.NearCastle, roadIdx, item.IsGroup)
		}
	}
}

func sectionContains(list []models.SidMapping, sid string) bool {
	for _, mapping := range list {
		if mapping.Sid == sid {
			return true
		}
	}
	return false
}

// collectZoneContentItems serialises every section into a flat list.
func (this *State) collectZoneContentItems() []models.ZoneContentItem {
	var out []models.ZoneContentItem
	collect := func(contentSection *zoneContentSection) {
		for _, row := range contentSection.rows {
			roadDistance := ""
			if row.RoadDistIdx >= 0 && row.RoadDistIdx < len(roadDistances) {
				roadDistance = roadDistances[row.RoadDistIdx]
			}
			out = append(out, models.ZoneContentItem{
				Sid:          row.Mapping.Sid,
				Name:         row.Mapping.Name,
				Count:        row.Count,
				IsGuarded:    row.IsGuarded.Value,
				NearCastle:   row.NearCastle.Value,
				RoadDistance: roadDistance,
				IsGroup:      row.IsGroup,
			})
		}
	}
	collect(this.zcMines)
	collect(this.zcTreasures)
	collect(this.zcHires)
	collect(this.zcBanks)
	return out
}
