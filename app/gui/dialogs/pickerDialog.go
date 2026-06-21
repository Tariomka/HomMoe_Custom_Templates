package dialogs

import (
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

// checkedRowBg is the highlight behind a selected picker row.
var checkedRowBg = themes.ColorSelection

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
	footer       func(gtx layout.Context, theme *material.Theme) layout.Dimensions
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
	rows := this.buildRows(theme, filter)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return widgets.NewTextboxWidget(theme, &this.search, "Search…")(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(rows) == 0 {
				label := material.Body2(theme, "(no matches)")
				label.Color = themes.ColorTextDim
				return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, label.Layout)
			}
			return material.List(theme, &this.scroll).Layout(gtx, len(rows), func(gtx layout.Context, index int) layout.Dimensions {
				return rows[index](gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if this.footer == nil {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return this.footer(gtx, theme)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutButtons(gtx, theme)
		}),
	), false
}

func (this *multiSelectPicker) layoutButtons(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	count := len(this.selectedIds())
	addText := this.addLabel
	if count > 1 {
		addText = this.addLabel + " (" + strconv.Itoa(count) + ")"
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 0)}
		}),
		layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelBtn, false)),
		layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
		layout.Rigid(widgets.NewGoldButtonWidget(theme, addText, &this.addBtn, count == 0)),
	)
}

// buildRows flattens the filtered entries into a list of group headers and leaf
// rows, preserving the entry order (so callers control grouping/sorting).
func (this *multiSelectPicker) buildRows(theme *material.Theme, filter string) []layout.Widget {
	var rows []layout.Widget
	emitted := map[string]bool{}

	appendGroup := func(group string) {
		// Count matching entries in this group for the header badge.
		count := 0
		for _, entry := range this.entries {
			if entry.group == group && this.matches(entry, filter) {
				count++
			}
		}
		rows = append(rows, this.groupHeader(theme, group, count))
	}

	for _, entry := range this.entries {
		if !this.matches(entry, filter) {
			continue
		}
		if this.grouped && !emitted[entry.group] {
			emitted[entry.group] = true
			appendGroup(entry.group)
		}
		rows = append(rows, this.leafRow(theme, entry))
	}
	return rows
}

func (this *multiSelectPicker) matches(entry pickerEntry, filter string) bool {
	return filter == "" || strings.Contains(entry.haystack, filter)
}

func (this *multiSelectPicker) groupHeader(theme *material.Theme, group string, count int) layout.Widget {
	col := themes.ColorAccent
	if this.groupColor != nil {
		col = this.groupColor(group)
	}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(theme, group)
					label.Color = col
					label.TextSize = unit.Sp(13)
					label.Font = font.Font{Weight: font.SemiBold}
					return label.Layout(gtx)
				}),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, "("+strconv.Itoa(count)+")")
					label.Color = themes.ColorTextDim
					label.TextSize = unit.Sp(11)
					return label.Layout(gtx)
				}),
			)
		})
	}
}

func (this *multiSelectPicker) leafRow(theme *material.Theme, entry pickerEntry) layout.Widget {
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
			dims := layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(16))
						mark := " "
						if checked {
							mark = "✓"
						}
						label := material.Body1(theme, mark)
						label.Color = themes.ColorAccentBright
						label.Font = font.Font{Weight: font.Bold}
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if entry.badge == "" {
							return layout.Dimensions{}
						}
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(34))
						label := material.Body2(theme, entry.badge)
						label.Color = themes.ColorTextDim
						label.TextSize = unit.Sp(11)
						return label.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(theme, entry.label)
						label.Color = themes.ColorText
						label.TextSize = unit.Sp(13)
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if entry.trailing == "" {
							return layout.Dimensions{}
						}
						label := material.Body2(theme, entry.trailing)
						label.Color = themes.ColorTextDim
						label.TextSize = unit.Sp(11)
						return label.Layout(gtx)
					}),
				)
			})
			call := macro.Stop()
			if checked {
				paint.FillShape(gtx.Ops, checkedRowBg, clip.Rect{Max: dims.Size}.Op())
			}
			call.Add(gtx.Ops)
			return dims
		})
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

// ── Item picker ─────────────────────────────────────────────────────────────

// NewItemPickerDialog builds a category-grouped artifact picker. excluded ids
// are hidden (already-banned / already-chosen). onApply receives the selected
// artifact SIDs.
func NewItemPickerDialog(title string, excluded []string, onApply func(ids []string)) *multiSelectPicker {
	skip := toSet(excluded)
	visible := make([]constants.BannableItemEntry, 0, len(constants.BannableItems))
	for _, item := range constants.BannableItems {
		if !skip[item.Sid] {
			visible = append(visible, item)
		}
	}
	sortStable(visible, func(a, b constants.BannableItemEntry) bool {
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		return a.Name < b.Name
	})

	entries := make([]pickerEntry, 0, len(visible))
	for _, item := range visible {
		entries = append(entries, pickerEntry{
			id:       item.Sid,
			group:    item.Category,
			label:    item.Name,
			trailing: item.Sid,
			haystack: strings.ToLower(item.Name + " " + item.Sid + " " + item.Category),
		})
	}

	picker := newMultiSelectPicker(title, entries, true)
	picker.onApply = onApply
	return picker
}

// ── Spell picker ────────────────────────────────────────────────────────────

// NewSpellPickerDialog builds a school-grouped, tier-sorted spell picker.
// excluded ids are hidden. When showMakeFree is true a "make free" toggle is
// shown. onApply receives the selected spell SIDs and the make-free flag.
func NewSpellPickerDialog(excluded []string, showMakeFree bool, onApply func(ids []string, makeFree bool)) *multiSelectPicker {
	skip := toSet(excluded)
	schoolRank := map[string]int{}
	for i, school := range constants.SpellSchoolOrder {
		schoolRank[school] = i
	}

	visible := make([]constants.SpellEntry, 0, len(constants.KnownSpells))
	for _, spell := range constants.KnownSpells {
		if !skip[spell.Sid] {
			visible = append(visible, spell)
		}
	}
	sortStable(visible, func(a, b constants.SpellEntry) bool {
		ra, rb := rankOf(schoolRank, a.School), rankOf(schoolRank, b.School)
		if ra != rb {
			return ra < rb
		}
		if a.Tier != b.Tier {
			return a.Tier < b.Tier
		}
		return a.Name < b.Name
	})

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
			badge:    "[T" + strconv.Itoa(spell.Tier) + "]",
			haystack: strings.ToLower(spell.Name + " " + spell.Sid + " " + spell.School),
		})
	}

	makeFree := new(widget.Bool)
	picker := newMultiSelectPicker("Pick Spells", entries, true)
	picker.groupColor = constants.GetSpellSchoolColorFromDisplayName
	if showMakeFree {
		picker.footer = func(gtx layout.Context, theme *material.Theme) layout.Dimensions {
			return widgets.NewLabeledCheckboxRowWidget(theme, makeFree, "Make spell(s) free")(gtx)
		}
	}
	picker.onApply = func(ids []string) {
		onApply(ids, makeFree.Value)
	}
	return picker
}

// ── Value-override picker ───────────────────────────────────────────────────

// NewValueOverridePickerDialog builds a flat SID picker with a single guard
// value applied to every selection. excluded SIDs are hidden. onApply receives
// "sid=guardValue" lines.
func NewValueOverridePickerDialog(excluded []string, onApply func(lines []string)) *multiSelectPicker {
	skip := toSet(excluded)
	sids := make([]string, 0, len(constants.ValueOverrideSids))
	for _, sid := range constants.ValueOverrideSids {
		if !skip[sid] {
			sids = append(sids, sid)
		}
	}
	sortStable(sids, func(a, b string) bool { return a < b })

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
	picker.footer = func(gtx layout.Context, theme *material.Theme) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(100))
				label := material.Body2(theme, "Guard value")
				label.Color = themes.ColorTextDim
				label.TextSize = unit.Sp(12)
				return label.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return widgets.NewTextboxWidget(theme, guardEdit, "5000")(gtx)
			}),
		)
	}
	picker.onApply = func(ids []string) {
		guard := 5000
		if parsed, err := strconv.Atoi(strings.TrimSpace(guardEdit.Text())); err == nil {
			guard = parsed
		}
		lines := make([]string, 0, len(ids))
		for _, sid := range ids {
			lines = append(lines, sid+"="+strconv.Itoa(guard))
		}
		onApply(lines)
	}
	return picker
}

// ── small helpers ───────────────────────────────────────────────────────────

func toSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

func rankOf(ranks map[string]int, key string) int {
	if rank, ok := ranks[key]; ok {
		return rank
	}
	return 99
}

// sortStable is a tiny insertion sort that keeps the original order of equal
// elements; the picker lists are small (<200) so this is more than adequate and
// avoids pulling generics gymnastics through the sort package.
func sortStable[T any](items []T, less func(a, b T) bool) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && less(items[j], items[j-1]); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}
