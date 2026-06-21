package drivers

import (
	"fmt"
	"iter"
	"sort"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
)

// zoneContentRow is one editable item inside a zone-content section. The legacy
// flat fields (guarded / near-castle / road-distance) are replaced by a
// polymorphic rule list edited through the Manage Rules dialog.
type zoneContentRow struct {
	Mapping models.SidMapping
	Count   int
	IsGroup bool
	rules   []models.ContentRuleRowSave

	countSld  widget.Float
	manageBtn widget.Clickable
	removeBtn widget.Clickable
	dupBtn    widget.Clickable
}

func newZoneContentRow(mapping models.SidMapping, count int, rules []models.ContentRuleRowSave, isGroup bool) *zoneContentRow {
	return &zoneContentRow{
		Mapping: mapping,
		Count:   count,
		IsGroup: isGroup,
		rules:   utils.CloneRuleRows(rules),
	}
}

// Rules returns a defensive copy of the row's content rules, letting the parent
// panel serialize them without exposing the row's mutable slice.
func (this *zoneContentRow) Rules() []models.ContentRuleRowSave {
	return utils.CloneRuleRows(this.rules)
}

// ZoneContentSection is one of the mandatory-content groups.
type ZoneContentSection struct {
	Title      string
	Items      []models.SidMapping
	MaxCount   int
	ShowNear   bool
	rows       []*zoneContentRow
	addPreset  *components.DropdownSelector
	addBtn     widget.Clickable
	openDialog interfaces.DialogOpener
}

func NewZoneContentSection(title string, items []models.SidMapping, maxCount int, showNear bool) *ZoneContentSection {
	// Present the "add content" dropdown alphabetically by display name. Sort a
	// copy so the shared ContentItemGroup global keeps its authored order.
	sorted := make([]models.SidMapping, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].Name) < strings.ToLower(sorted[j].Name)
	})
	labels := make([]string, len(sorted))
	for i, item := range sorted {
		labels[i] = item.Name
	}
	return &ZoneContentSection{
		Title:     title,
		Items:     sorted,
		MaxCount:  maxCount,
		ShowNear:  showNear,
		addPreset: components.NewDropdownSelector(labels),
	}
}

// SetDialogOpener wires the section to the modal host so rows can launch the
// Manage Rules dialog.
func (this *ZoneContentSection) SetDialogOpener(opener interfaces.DialogOpener) {
	this.openDialog = opener
}

// Add appends a new row using the given mapping and rule list.
func (this *ZoneContentSection) Add(mapping models.SidMapping, count int, rules []models.ContentRuleRowSave, group bool) {
	if count < 1 {
		count = 1
	}
	if count > this.MaxCount {
		count = this.MaxCount
	}
	this.rows = append(this.rows, newZoneContentRow(mapping, count, rules, group))
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
			this.Add(this.Items[idx], 1, defaultContentRules(), false)
		}
		// Process per-row clicks (collect indices to remove).
		keep := this.rows[:0]
		for i, row := range this.rows {
			if row.removeBtn.Clicked(gtx) {
				continue
			}
			if row.dupBtn.Clicked(gtx) {
				keep = append(keep, row)
				clone := newZoneContentRow(row.Mapping, row.Count, row.rules, row.IsGroup)
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
						label := material.Caption(theme, "Add preset:")
						label.Color = themes.ColorTextDim
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
					label := material.Caption(theme, "(no items)")
					label.Color = themes.ColorTextDim
					return layout.Inset{Top: unit.Dp(4), Left: unit.Dp(4)}.Layout(gtx, label.Layout)
				}
				children := make([]layout.FlexChild, 0, len(this.rows)*2)
				for i, row := range this.rows {
					if i > 0 {
						children = append(children, layout.Rigid(widgets.NewVerticalSpacerWidget(4)))
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

		// Launch the Manage Rules dialog. The captured row pointer keeps the
		// callback bound to this exact row across frames/reorders.
		if row.manageBtn.Clicked(gtx) && this.openDialog != nil {
			captured := row
			this.openDialog(dialogs.NewManageRulesDialog(captured.Mapping, captured.rules, func(updated []models.ContentRuleRowSave) {
				captured.rules = updated
			}))
		}

		return widgets.NewPanelWidget(unit.Dp(6), func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(theme, rowDisplayName(row))
							label.Color = themes.ColorAccent
							label.TextSize = unit.Sp(13)
							return label.Layout(gtx)
						}),
						layout.Rigid(widgets.NewButtonWidget(theme, "Duplicate", &row.dupBtn, false)),
						layout.Rigid(widgets.NewHorizontalSpacerWidget(4)),
						layout.Rigid(widgets.NewButtonWidget(theme, "Remove", &row.removeBtn, false)),
					)
				}),
				layout.Rigid(widgets.NewVerticalSpacerWidget(4)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(0.6, widgets.NewLabeledRowWidget(theme, "Count", 60, widgets.NewLabeledSliderWidget(theme, &row.countSld, fmt.Sprintf("%d", liveCount)))),
						layout.Rigid(widgets.NewHorizontalSpacerWidget(16)),
						layout.Flexed(0.4, widgets.NewLabeledRowWidget(theme, "Rules", 50, this.layoutMarkers(theme, row))),
						layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
						layout.Rigid(widgets.NewButtonWidget(theme, "Manage Rules", &row.manageBtn, false)),
					)
				}),
			)
		})(gtx)
	}
}

// layoutMarkers renders the compact marker badges for a row's rules, or a dim
// placeholder when the row has no rules.
func (this *ZoneContentSection) layoutMarkers(theme *material.Theme, row *zoneContentRow) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		markers := ruleMarkers(row.Mapping, row.rules)
		if markers == "" {
			label := material.Body2(theme, "(none)")
			label.Color = themes.ColorTextDim
			label.TextSize = unit.Sp(12)
			return label.Layout(gtx)
		}
		label := material.Body1(theme, markers)
		label.Color = themes.ColorAccentBright
		label.TextSize = unit.Sp(13)
		return label.Layout(gtx)
	}
}

// rowDisplayName mirrors the C# ZoneContentItemUI.DisplayName: the item name,
// with the chosen variant's description appended when a Variant rule applies.
func rowDisplayName(row *zoneContentRow) string {
	for _, saved := range row.rules {
		if saved.VariantId == nil || !strings.EqualFold(saved.Name, content_rules.RuleVariantName) {
			continue
		}
		for _, variant := range content_rules.GetVariantsForContent(row.Mapping) {
			if description, ok := variant.Variants[*saved.VariantId]; ok {
				return row.Mapping.Name + " (" + description + ")"
			}
		}
	}
	return row.Mapping.Name
}

// defaultContentRules is the rule list applied to a freshly-added row: a single
// Guarded rule, matching the historical default of the Guarded checkbox.
func defaultContentRules() []models.ContentRuleRowSave {
	guarded := true
	return []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &guarded}}
}

// ruleMarkers returns the concatenated marker badges for a rule list (e.g.
// "G · R · S"), skipping rules that have no marker (Variant) or are invalid.
func ruleMarkers(mapping models.SidMapping, rules []models.ContentRuleRowSave) string {
	parts := make([]string, 0, len(rules))
	for _, saved := range rules {
		rule := content_rules.CreateRuleFromSavedRule(saved, mapping)
		if rule == nil {
			continue
		}
		if marker := rule.Marker(); marker != "" {
			parts = append(parts, marker)
		}
	}
	return strings.Join(parts, " · ")
}
