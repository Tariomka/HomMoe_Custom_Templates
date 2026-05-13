package gui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
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

func newZoneContentRow(m models.SidMapping, count int, guarded, nearCastle bool, roadIdx int, isGroup bool) *zoneContentRow {
	r := &zoneContentRow{
		Mapping:     m,
		Count:       count,
		RoadDistIdx: roadIdx,
		IsGroup:     isGroup,
		roadCombo:   newComboBox(roadDistances),
	}
	if r.RoadDistIdx >= 0 && r.RoadDistIdx < len(roadDistances) {
		r.roadCombo.Selected = r.RoadDistIdx
	}
	r.IsGuarded.Value = guarded
	r.NearCastle.Value = nearCastle
	return r
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
	for i, it := range items {
		labels[i] = it.Name
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
func (sec *zoneContentSection) Add(m models.SidMapping, count int, guarded, near bool, roadIdx int, group bool) {
	if count < 1 {
		count = 1
	}
	if count > sec.MaxCount {
		count = sec.MaxCount
	}
	sec.rows = append(sec.rows, newZoneContentRow(m, count, guarded, near, roadIdx, group))
}

func (sec *zoneContentSection) Layout(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Process per-section button clicks first.
		if sec.addBtn.Clicked(gtx) && len(sec.Items) > 0 {
			idx := sec.addPreset.Selected
			if idx < 0 || idx >= len(sec.Items) {
				idx = 0
			}
			sec.Add(sec.Items[idx], 1, true, false, 0, false)
		}
		// Process per-row clicks (collect indices to remove).
		keep := sec.rows[:0]
		for i, r := range sec.rows {
			if r.removeBtn.Clicked(gtx) {
				continue
			}
			if r.dupBtn.Clicked(gtx) {
				keep = append(keep, r)
				clone := newZoneContentRow(r.Mapping, r.Count, r.IsGuarded.Value, r.NearCastle.Value, r.RoadDistIdx, r.IsGroup)
				keep = append(keep, clone)
				continue
			}
			_ = i
			keep = append(keep, r)
		}
		sec.rows = keep

		return section(th, sec.Title, []layout.Widget{
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(120)
						lbl := material.Body2(th, "Add preset:")
						lbl.Color = colTextDim
						lbl.TextSize = unit.Sp(12)
						return lbl.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return sec.addPreset.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return toolbarButton{Text: "+ Add", Click: &sec.addBtn}.Layout(gtx, th)
					}),
				)
			},
			func(gtx layout.Context) layout.Dimensions {
				if len(sec.rows) == 0 {
					lbl := material.Body2(th, "(no items)")
					lbl.Color = colTextDim
					lbl.TextSize = unit.Sp(12)
					return layout.Inset{Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, lbl.Layout)
				}
				children := make([]layout.FlexChild, 0, len(sec.rows)*2)
				for i, r := range sec.rows {
					r := r
					if i > 0 {
						children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
					}
					children = append(children, layout.Rigid(sec.layoutRow(th, r)))
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
			},
		})(gtx)
	}
}

func (sec *zoneContentSection) layoutRow(th *material.Theme, r *zoneContentRow) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Sync count slider into integer field.
		desired := mapRangeInv(float32(r.Count), 1, float32(sec.MaxCount))
		if !r.countSld.Dragging() && r.countSld.Value == 0 && r.Count > 0 {
			r.countSld.Value = desired
		}
		liveCount := roundedRange(r.countSld.Value, 1, sec.MaxCount)
		r.Count = liveCount
		r.RoadDistIdx = r.roadCombo.Selected

		return borderedPanel(gtx, unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(th, r.Mapping.Name)
							lbl.Color = colGold
							lbl.TextSize = unit.Sp(13)
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return toolbarButton{Text: "Duplicate", Click: &r.dupBtn}.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return toolbarButton{Text: "Remove", Click: &r.removeBtn}.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
				layout.Rigid(labeledRowW(th, "Count", 100, func(gtx layout.Context) layout.Dimensions {
					return sliderLabeled(gtx, th, &r.countSld, fmt.Sprintf("%d", liveCount))
				})),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(checkRow(th, &r.IsGuarded, "Guarded")),
						layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !sec.ShowNear {
								return layout.Dimensions{}
							}
							return checkRow(th, &r.NearCastle, "Near castle")(gtx)
						}),
					)
				}),
				layout.Rigid(labeledRowW(th, "Road distance", 100, func(gtx layout.Context) layout.Dimensions {
					return r.roadCombo.Layout(gtx, th)
				})),
			)
		})
	}
}

// — State integration —

// seedDefaultPlayerZoneContent mirrors C# InitializeDefaultPlayerZoneContents.
func (s *State) seedDefaultPlayerZoneContent() {
	s.zcMines.rows = nil
	s.zcTreasures.rows = nil
	s.zcHires.rows = nil
	s.zcBanks.rows = nil

	// Mines: wood/ore/gold guarded next-to-castle; crystals/mercury/gemstones/alchemy-lab guarded near road.
	s.zcMines.Add(constants.ContentIds.MineWood, 1, true, true, 0, false)
	s.zcMines.Add(constants.ContentIds.MineOre, 1, true, true, 0, false)
	s.zcMines.Add(constants.ContentIds.MineGold, 1, true, true, 0, false)
	s.zcMines.Add(constants.ContentIds.MineCrystals, 1, true, false, 1, false)
	s.zcMines.Add(constants.ContentIds.MineMercury, 1, true, false, 1, false)
	s.zcMines.Add(constants.ContentIds.MineGemstones, 1, true, false, 1, false)
	s.zcMines.Add(constants.ContentIds.AlchemyLab, 1, true, false, 1, false)

	// Treasures: PandoraBox + RandomItemEpic guarded.
	s.zcTreasures.Add(constants.ContentIds.PandoraBox, 1, true, false, 0, false)
	s.zcTreasures.Add(constants.ContentIds.RandomItemEpic, 1, true, false, 0, false)

	// Random hires: low ×2, high ×1, all-tier ×1 (groups).
	s.zcHires.Add(constants.IncludeListIds.RandomHiresLowTier, 2, true, false, 0, true)
	s.zcHires.Add(constants.IncludeListIds.RandomHiresHighTier, 1, true, false, 0, true)
	s.zcHires.Add(constants.IncludeListIds.RandomHiresAllTier, 1, true, false, 0, true)

	// Resource banks: tier1 ×2, tier2 ×1.
	s.zcBanks.Add(constants.IncludeListIds.ResourceBanksTier1, 2, true, false, 0, true)
	s.zcBanks.Add(constants.IncludeListIds.ResourceBanksTier2, 1, true, false, 0, true)
}

// applyZoneContentItems replaces every section based on a flat list of items
// loaded from a settings file. Items are routed to the appropriate section by
// SID lookup.
func (s *State) applyZoneContentItems(items []models.ZoneContentItem) {
	s.zcMines.rows = nil
	s.zcTreasures.rows = nil
	s.zcHires.rows = nil
	s.zcBanks.rows = nil
	for _, it := range items {
		m := models.SidMapping{Sid: it.Sid, Name: it.Name}
		if found, ok := helpers.LookupSid(it.Sid); ok {
			m = found
		}
		count := it.Count
		if count < 1 {
			count = 1
		}
		roadIdx := indexOf(roadDistances, it.RoadDistance)
		if roadIdx < 0 {
			roadIdx = 0
		}
		switch {
		case sectionContains(constants.ContentItemGroup.Mines, it.Sid):
			s.zcMines.Add(m, count, it.IsGuarded, it.NearCastle, roadIdx, it.IsGroup)
		case sectionContains(constants.ContentItemGroup.Treasures, it.Sid):
			s.zcTreasures.Add(m, count, it.IsGuarded, it.NearCastle, roadIdx, it.IsGroup)
		case sectionContains(constants.ContentItemGroup.HireBuildings, it.Sid):
			s.zcHires.Add(m, count, it.IsGuarded, it.NearCastle, roadIdx, it.IsGroup)
		case sectionContains(constants.ContentItemGroup.ResourceBanks, it.Sid):
			s.zcBanks.Add(m, count, it.IsGuarded, it.NearCastle, roadIdx, it.IsGroup)
		default:
			// Unknown SID — keep with treasures by default.
			s.zcTreasures.Add(m, count, it.IsGuarded, it.NearCastle, roadIdx, it.IsGroup)
		}
	}
}

func sectionContains(list []models.SidMapping, sid string) bool {
	for _, m := range list {
		if m.Sid == sid {
			return true
		}
	}
	return false
}

// collectZoneContentItems serialises every section into a flat list.
func (s *State) collectZoneContentItems() []models.ZoneContentItem {
	var out []models.ZoneContentItem
	collect := func(sec *zoneContentSection) {
		for _, r := range sec.rows {
			rd := ""
			if r.RoadDistIdx >= 0 && r.RoadDistIdx < len(roadDistances) {
				rd = roadDistances[r.RoadDistIdx]
			}
			out = append(out, models.ZoneContentItem{
				Sid:          r.Mapping.Sid,
				Name:         r.Mapping.Name,
				Count:        r.Count,
				IsGuarded:    r.IsGuarded.Value,
				NearCastle:   r.NearCastle.Value,
				RoadDistance: rd,
				IsGroup:      r.IsGroup,
			})
		}
	}
	collect(s.zcMines)
	collect(s.zcTreasures)
	collect(s.zcHires)
	collect(s.zcBanks)
	return out
}
