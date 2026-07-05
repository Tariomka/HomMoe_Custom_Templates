package components

import (
	"image"

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

type dropdownItem struct {
	label string
	row   widget.Clickable
}

func newDropdownItem(label string) *dropdownItem {
	return &dropdownItem{label: label}
}

type DropdownSelector struct {
	items         []*dropdownItem
	selectedIndex int

	toggle widget.Clickable
	isOpen bool

	WasUpdated bool
}

func NewDropdownSelector(labels []string) *DropdownSelector {
	items := make([]*dropdownItem, len(labels))
	for i, label := range labels {
		items[i] = newDropdownItem(label)
	}
	return &DropdownSelector{items: items}
}

func (this *DropdownSelector) GetSelectedIndex() int {
	return this.selectedIndex
}

// SetItems replaces the option list and resets selection bounds.
func (this *DropdownSelector) SetItems(items []string) {
	newItems := make([]*dropdownItem, len(items))
	for i, label := range items {
		newItems[i] = newDropdownItem(label)
	}
	this.items = newItems
	if this.selectedIndex >= len(this.items) {
		this.selectedIndex = 0
	}
}

// SelectByName sets Selected to the index whose label matches name.
func (this *DropdownSelector) SelectByName(name string) bool {
	for i, item := range this.items {
		if item.label == name {
			this.selectedIndex = i
			return true
		}
	}
	return false
}

func (this *DropdownSelector) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		this.update(gtx)
		children := []layout.FlexChild{layout.Rigid(this.getTriggerWidget(theme))}
		if this.isOpen && len(this.items) > 0 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, this.getListWidget(theme))
			}))
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	}
}

func (this *DropdownSelector) getTriggerWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		arrowText := "▼"
		borderColor := themes.ColorBorder
		if this.toggle.Hovered() {
			borderColor = themes.ColorHover
		}
		if this.isOpen {
			arrowText = "▲"
			borderColor = themes.ColorAccent
		}
		return material.Clickable(gtx, &this.toggle, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.UniformInset(constants.DefaultPaddingSmall).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, widgets.NewLabelBuilder(theme).
							WithSizeDefault().WithText(this.value()).WithColorDefault().WithMaxLines(1).Build),
						layout.Rigid(widgets.NewLabelBuilder(theme).
							WithSize(10).WithText(arrowText).WithColor(themes.ColorAccentDim).Build))
				})
			call := macro.Stop()
			radius := gtx.Dp(constants.DefaultRoundness)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, themes.ColorInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func (this *DropdownSelector) getListWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		rows := []layout.FlexChild{}
		for i, item := range this.items {
			rows = append(rows, layout.Rigid(
				widgets.NewDropdownRowButtonWidget(theme, item.label, &item.row, i == this.selectedIndex)))
		}
		dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
		call := macro.Stop()
		radius := gtx.Dp(constants.DefaultRoundness)
		rect := image.Rectangle{Max: dims.Size}
		paint.FillShape(gtx.Ops, themes.ColorInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, themes.ColorBorder, clip.Stroke{
			Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
			Width: float32(gtx.Dp(1)),
		}.Op())
		call.Add(gtx.Ops)
		return dims
	}
}

// update returns true if the selection changed this frame.
func (this *DropdownSelector) update(gtx layout.Context) bool {
	this.WasUpdated = false
	if this.toggle.Clicked(gtx) {
		this.isOpen = !this.isOpen
	}
	for i, item := range this.items {
		if item.row.Clicked(gtx) {
			if this.selectedIndex != i {
				this.selectedIndex = i
				this.WasUpdated = true
			}
			this.isOpen = false
		}
	}
	return this.WasUpdated
}

// value returns the currently selected option label, or "" if empty.
func (this *DropdownSelector) value() string {
	if this.selectedIndex >= 0 && this.selectedIndex < len(this.items) {
		return this.items[this.selectedIndex].label
	}

	return ""
}
