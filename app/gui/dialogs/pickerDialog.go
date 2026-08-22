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
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
)

// MultiSelectPicker is a generic searchable, optionally-grouped, multi-select
// modal. It backs the item, spell and value-override pickers; per-picker extras
// (a guard-value field, a "make free" toggle) are supplied via the footer hook
// and read back inside onApply. Implements widgets.Dialog.
type MultiSelectPicker struct {
	title         string
	prefW, prefH  unit.Dp
	entries       []dtos.PickerEntryDto
	grouped       bool
	groupColor    func(group string) color.NRGBA
	addLabel      string
	footerWidget  func(theme *material.Theme) layout.Widget
	onApply       func(ids []string)
	pickerHandler handler_interfaces.IPickerHandler

	search    widget.Editor
	scroll    widget.List
	selected  map[string]bool
	clicks    map[string]*widget.Clickable
	addBtn    widget.Clickable
	cancelBtn widget.Clickable
}

func NewItemPickerDialog(
	title string,
	excluded []string,
	pickerHandler handler_interfaces.IPickerHandler,
	onApply func(ids []string)) *MultiSelectPicker {
	visible := constants.GetBannableItemsWithExclusions(excluded)

	items := make([]dtos.PickerItemDto, 0, len(visible))
	for _, item := range visible {
		items = append(items, dtos.PickerItemDto{
			Sid:      item.Sid,
			Name:     item.Name,
			Category: item.Category,
		})
	}

	picker := newMultiSelectPicker(
		title, pickerHandler.BuildItemPickerEntries(items), true, pickerHandler)
	picker.onApply = onApply
	return picker
}

func NewSpellPickerDialog(
	excluded []string,
	showMakeFree bool,
	pickerHandler handler_interfaces.IPickerHandler,
	onApply func(ids []string, makeFree bool)) *MultiSelectPicker {
	visible := constants.GetKnownSpellsWithExclusions(excluded)

	spells := make([]dtos.PickerSpellDto, 0, len(visible))
	for _, spell := range visible {
		spells = append(spells, dtos.PickerSpellDto{
			Sid:               spell.Sid,
			Name:              spell.Name,
			School:            spell.School,
			SchoolDisplayName: constants.GetSpellSchoolDisplayName(spell.School),
			Tier:              spell.Tier,
		})
	}

	makeFree := new(widget.Bool)
	picker := newMultiSelectPicker(
		"Pick Spells", pickerHandler.BuildSpellPickerEntries(spells), true, pickerHandler)
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

func NewValueOverridePickerDialog(
	excluded []string,
	pickerHandler handler_interfaces.IPickerHandler,
	onApply func(lines []string)) *MultiSelectPicker {
	sids := constants.GetValueOverrideSidsWithExclusions(excluded)

	guardEdit := &widget.Editor{SingleLine: true}
	guardEdit.SetText("5000")
	picker := newMultiSelectPicker(
		"Pick Value Overrides",
		pickerHandler.BuildValueOverridePickerEntries(sids),
		false,
		pickerHandler)
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

func newMultiSelectPicker(
	title string,
	entries []dtos.PickerEntryDto,
	grouped bool,
	pickerHandler handler_interfaces.IPickerHandler) *MultiSelectPicker {
	picker := &MultiSelectPicker{
		title:         title,
		prefW:         unit.Dp(560),
		prefH:         unit.Dp(560),
		entries:       entries,
		grouped:       grouped,
		addLabel:      "Add Selected",
		selected:      map[string]bool{},
		clicks:        map[string]*widget.Clickable{},
		pickerHandler: pickerHandler,
	}
	picker.search.SingleLine = true
	picker.scroll.Axis = layout.Vertical
	return picker
}

func (this *MultiSelectPicker) Title() string { return this.title }

func (this *MultiSelectPicker) PreferredSize() (unit.Dp, unit.Dp) { return this.prefW, this.prefH }

func (this *MultiSelectPicker) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	if this.addBtn.Clicked(gtx) {
		if this.onApply != nil {
			this.onApply(this.selectedIDs())
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	if this.cancelBtn.Clicked(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	filter := this.pickerHandler.NormalizePickerFilter(this.search.Text())
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

func (this *MultiSelectPicker) getButtonsWidget(theme *material.Theme) layout.Widget {
	count := len(this.selectedIDs())
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

// getRowWidgets renders the filtered row model produced by the picker handler.
func (this *MultiSelectPicker) getRowWidgets(theme *material.Theme, filter string) []layout.Widget {
	model := this.pickerHandler.GetVisiblePickerRows(this.entries, filter, this.grouped)

	rows := make([]layout.Widget, 0, len(model))
	for _, row := range model {
		if row.IsGroupHeader {
			rows = append(rows, this.getGroupHeaderWidget(theme, row.Group, row.GroupMatchCount))

			continue
		}
		rows = append(rows, this.getLeafRowWidget(theme, row.Entry))
	}

	return rows
}

func (this *MultiSelectPicker) getGroupHeaderWidget(theme *material.Theme, group string, count int) layout.Widget {
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

func (this *MultiSelectPicker) getLeafRowWidget(theme *material.Theme, entry dtos.PickerEntryDto) layout.Widget {
	clk := this.clickFor(entry.ID)
	return func(gtx layout.Context) layout.Dimensions {
		if clk.Clicked(gtx) {
			if this.selected[entry.ID] {
				delete(this.selected, entry.ID)
			} else {
				this.selected[entry.ID] = true
			}
		}
		checked := this.selected[entry.ID]
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
func getLeafRowContentWidget(theme *material.Theme, entry dtos.PickerEntryDto, checked bool) layout.Widget {
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
				if entry.Badge == "" {
					return layout.Dimensions{}
				}

				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(34))
				label := material.Overline(theme, entry.Badge)
				label.Color = themes.ColorsBase.TextDim
				return label.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, entry.Label)
				label.Color = themes.ColorsBase.Text
				return label.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if entry.Trailing == "" {
					return layout.Dimensions{}
				}

				label := material.Overline(theme, entry.Trailing)
				label.Color = themes.ColorsBase.TextDim
				return label.Layout(gtx)
			}),
		)
	}
}

func (this *MultiSelectPicker) clickFor(id string) *widget.Clickable {
	clk := this.clicks[id]
	if clk == nil {
		clk = &widget.Clickable{}
		this.clicks[id] = clk
	}
	return clk
}

func (this *MultiSelectPicker) selectedIDs() []string {
	return this.pickerHandler.GetSelectedPickerIDs(this.entries, this.selected)
}
