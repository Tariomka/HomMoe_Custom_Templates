package dialogs

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

// pickerEntry is one selectable row in a multiSelectPicker.
type pickerEntry struct {
	id       string
	group    string // category / school; ignored when the picker is flat
	label    string // primary display text
	badge    string // optional leading badge (e.g. "[T3]")
	trailing string // optional dim trailing text (e.g. the raw SID)
	haystack string // lowercased search text
}

// multiSelectPicker is a generic searchable, optionally-grouped, multi-select
// modal. It backs the item, spell and value-override pickers; per-picker extras
// (a guard-value field, a "make free" toggle) are supplied via the footer hook
// and read back inside onApply. Implements widgets.Dialog.
type multiSelectPicker struct {
	title        string
	prefW, prefH unit.Dp
	entries      []pickerEntry
	grouped      bool
	groupColor   func(group string) color.NRGBA
	addLabel     string
	footerWidget func(theme *material.Theme) layout.Widget
	onApply      func(ids []string)

	search    widget.Editor
	scroll    widget.List
	selected  map[string]bool
	clicks    map[string]*widget.Clickable
	addBtn    widget.Clickable
	cancelBtn widget.Clickable
}

func newMultiSelectPicker(title string, entries []pickerEntry, grouped bool) *multiSelectPicker {
	picker := &multiSelectPicker{
		title:    title,
		prefW:    unit.Dp(560),
		prefH:    unit.Dp(560),
		entries:  entries,
		grouped:  grouped,
		addLabel: "Add Selected",
		selected: map[string]bool{},
		clicks:   map[string]*widget.Clickable{},
	}
	picker.search.SingleLine = true
	picker.scroll.Axis = layout.Vertical
	return picker
}

func NewItemPickerDialog(title string, excluded []string, onApply func(ids []string)) *multiSelectPicker {
	visible := constants.GetBannableItemsWithExclusions(excluded)

	entries := make([]pickerEntry, 0, len(visible))
	for _, item := range visible {
		entries = append(entries, pickerEntry{
			id:       item.Sid,
			group:    item.Category,
			label:    item.Name,
			haystack: strings.ToLower(item.Name + " " + item.Sid + " " + item.Category),
		})
	}

	picker := newMultiSelectPicker(title, entries, true)
	picker.onApply = onApply
	return picker
}

func NewSpellPickerDialog(
	excluded []string,
	showMakeFree bool,
	onApply func(ids []string, makeFree bool)) *multiSelectPicker {
	visible := constants.GetKnownSpellsWithExclusions(excluded)

	entries := make([]pickerEntry, 0, len(visible))
	for _, spell := range visible {
		label := constants.SpellSchoolDisplayNames[spell.School]
		if label == "" {
			label = spell.School
		}
		entries = append(entries, pickerEntry{
			id:       spell.Sid,
			group:    label,
			label:    spell.Name,
			badge:    fmt.Sprintf("[T%d]", spell.Tier),
			haystack: strings.ToLower(spell.Name + " " + spell.Sid + " " + spell.School),
		})
	}

	makeFree := new(widget.Bool)
	picker := newMultiSelectPicker("Pick Spells", entries, true)
	picker.groupColor = constants.GetSpellSchoolColorFromDisplayName
	if showMakeFree {
		picker.footerWidget = func(theme *material.Theme) layout.Widget {
			return widgets.NewLabeledCheckboxRowWidget(theme, makeFree, "Make spell(s) free")
		}
	}
	picker.onApply = func(ids []string) {
		onApply(ids, makeFree.Value)
	}
	return picker
}

func NewValueOverridePickerDialog(excluded []string, onApply func(lines []string)) *multiSelectPicker {
	sids := constants.GetValueOverrideSidsWithExclusions(excluded)

	entries := make([]pickerEntry, 0, len(sids))
	for _, sid := range sids {
		entries = append(entries, pickerEntry{
			id:       sid,
			label:    sid,
			haystack: strings.ToLower(sid),
		})
	}

	guardEdit := &widget.Editor{SingleLine: true}
	guardEdit.SetText("5000")
	picker := newMultiSelectPicker("Pick Value Overrides", entries, false)
	picker.footerWidget = func(theme *material.Theme) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(100))
					label := material.Caption(theme, "Guard value")
					label.Color = themes.ColorsBase.TextDim
					return label.Layout(gtx)
				}),
				layout.Flexed(1, widgets.NewTextboxWidget(theme, guardEdit, "5000", false)))
		}
	}
	picker.onApply = func(ids []string) {
		guard := 5000
		if parsed, err := strconv.Atoi(strings.TrimSpace(guardEdit.Text())); err == nil {
			guard = parsed
		}
		lines := make([]string, 0, len(ids))
		for _, sid := range ids {
			lines = append(lines, fmt.Sprintf("%s=%d", sid, guard))
		}
		onApply(lines)
	}
	return picker
}

func (this *multiSelectPicker) Title() string { return this.title }

func (this *multiSelectPicker) PreferredSize() (unit.Dp, unit.Dp) { return this.prefW, this.prefH }

func (this *multiSelectPicker) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.addBtn.Clicked(gtx) {
		if this.onApply != nil {
			this.onApply(this.selectedIds())
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	if this.cancelBtn.Clicked(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	filter := strings.ToLower(strings.TrimSpace(this.search.Text()))
	rows := this.getRowWidgets(theme, filter)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(widgets.NewTextboxWidget(theme, &this.search, "Search...", false)),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(rows) == 0 {
				label := material.Body2(theme, "(no matches)")
				label.Color = themes.ColorsBase.TextDim
				return layout.Inset{Top: constants.DefaultPaddingSmall - 2}.Layout(gtx, label.Layout)
			}

			return material.List(theme, &this.scroll).Layout(gtx, len(rows),
				func(gtx layout.Context, index int) layout.Dimensions { return rows[index](gtx) })
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if this.footerWidget == nil {
				return layout.Dimensions{}
			}

			return layout.Inset{Top: constants.DefaultPadding}.Layout(gtx, this.footerWidget(theme))
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(10)),
		layout.Rigid(this.getButtonsWidget(theme)),
	), false
}

func (this *multiSelectPicker) getButtonsWidget(theme *material.Theme) layout.Widget {
	count := len(this.selectedIds())
	addText := this.addLabel
	if count > 1 {
		addText = fmt.Sprintf("%s (%d)", this.addLabel, count)
	}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 0)}
			}),
			layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelBtn, false)),
			widgets.NewDefaultComponentSpacer(),
			layout.Rigid(widgets.NewBrightButtonWidget(theme, addText, &this.addBtn, count == 0)),
		)
	}
}

// getRowWidgets flattens the filtered entries into a list of group headers and leaf
// rows, preserving the entry order (so callers control grouping/sorting).
func (this *multiSelectPicker) getRowWidgets(theme *material.Theme, filter string) []layout.Widget {
	var rows []layout.Widget
	emitted := map[string]bool{}

	appendGroup := func(group string) {
		// Count matching entries in this group for the header badge.
		count := 0
		for _, entry := range this.entries {
			if entry.group == group && strings.Contains(entry.haystack, filter) {
				count++
			}
		}
		rows = append(rows, this.getGroupHeaderWidget(theme, group, count))
	}

	for _, entry := range this.entries {
		if !strings.Contains(entry.haystack, filter) {
			continue
		}
		if this.grouped && !emitted[entry.group] {
			emitted[entry.group] = true
			appendGroup(entry.group)
		}
		rows = append(rows, this.getLeafRowWidget(theme, entry))
	}
	return rows
}

func (this *multiSelectPicker) getGroupHeaderWidget(theme *material.Theme, group string, count int) layout.Widget {
	groupTextColor := themes.ColorsBase.Accent
	if this.groupColor != nil {
		groupTextColor = this.groupColor(group)
	}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPaddingSmall, Bottom: constants.DefaultPaddingSmall - 4}.
			Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(theme, group)
						label.Color = groupTextColor
						label.Font = font.Font{Weight: font.SemiBold}
						return label.Layout(gtx)
					}),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Overline(theme, fmt.Sprintf("(%d)", count))
						label.Color = themes.ColorsBase.TextDim
						return label.Layout(gtx)
					}),
				)
			})
	}
}

func (this *multiSelectPicker) getLeafRowWidget(theme *material.Theme, entry pickerEntry) layout.Widget {
	clk := this.clickFor(entry.id)
	return func(gtx layout.Context) layout.Dimensions {
		if clk.Clicked(gtx) {
			if this.selected[entry.id] {
				delete(this.selected, entry.id)
			} else {
				this.selected[entry.id] = true
			}
		}
		checked := this.selected[entry.id]
		return material.Clickable(gtx, clk, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(constants.DefaultPaddingSmall-2).
				Layout(gtx, getLeafRowContentWidget(theme, entry, checked))
			call := macro.Stop()
			if checked {
				paint.FillShape(gtx.Ops, themes.ColorsBase.Selection, clip.Rect{Max: dims.Size}.Op())
			} else if clk.Hovered() {
				paint.FillShape(gtx.Ops, themes.ColorsBase.Hover, clip.Rect{Max: dims.Size}.Op())
			}
			call.Add(gtx.Ops)
			return dims
		})
	}
}

// getLeafRowContentWidget lays out a leaf row's checkmark, optional badge,
// label and optional trailing text.
func getLeafRowContentWidget(theme *material.Theme, entry pickerEntry, checked bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(16))
				mark := " "
				if checked {
					mark = "v"
				}
				label := material.Body1(theme, mark)
				label.Color = themes.ColorsBase.AccentBright
				label.Font = font.Font{Weight: font.Bold, Style: font.Italic}
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if entry.badge == "" {
					return layout.Dimensions{}
				}

				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(34))
				label := material.Overline(theme, entry.badge)
				label.Color = themes.ColorsBase.TextDim
				return label.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, entry.label)
				label.Color = themes.ColorsBase.Text
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if entry.trailing == "" {
					return layout.Dimensions{}
				}

				label := material.Overline(theme, entry.trailing)
				label.Color = themes.ColorsBase.TextDim
				return label.Layout(gtx)
			}),
		)
	}
}

func (this *multiSelectPicker) clickFor(id string) *widget.Clickable {
	clk := this.clicks[id]
	if clk == nil {
		clk = &widget.Clickable{}
		this.clicks[id] = clk
	}
	return clk
}

func (this *multiSelectPicker) selectedIds() []string {
	var ids []string
	for _, entry := range this.entries {
		if this.selected[entry.id] {
			ids = append(ids, entry.id)
		}
	}
	return ids
}
