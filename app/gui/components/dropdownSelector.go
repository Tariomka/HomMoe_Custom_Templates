package components

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
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

func (this *DropdownSelector) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.update(gtx)
	flex := layout.Flex{Axis: layout.Vertical}
	children := []layout.FlexChild{
		layout.Rigid(this.getTriggerWidget(theme)),
	}
	if this.isOpen && len(this.items) > 0 {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return this.layoutList(gtx, theme)
			})
		}))
	}
	return flex.Layout(gtx, children...)
}

func (this *DropdownSelector) getTriggerWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return material.Clickable(gtx, &this.toggle, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(theme, this.value())
						label.Color = themes.ColorText
						label.TextSize = unit.Sp(13)
						label.MaxLines = 1
						label.Truncator = "…"
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						arrow := "▾"
						if this.isOpen {
							arrow = "▴"
						}
						label := material.Body1(theme, arrow)
						label.Color = themes.ColorAccentDim
						label.TextSize = unit.Sp(11)
						return label.Layout(gtx)
					}),
				)
			})
			call := macro.Stop()
			radius := gtx.Dp(2)
			rect := image.Rectangle{Max: dims.Size}
			paint.FillShape(gtx.Ops, themes.ColorInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
			border := themes.ColorBorder
			if this.isOpen {
				border = themes.ColorAccent
			}
			paint.FillShape(gtx.Ops, border, clip.Stroke{
				Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
				Width: float32(gtx.Dp(1)),
			}.Op())
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func (this *DropdownSelector) layoutList(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	rows := make([]layout.FlexChild, 0, len(this.items))
	for i, item := range this.items {
		rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Clickable(gtx, &item.row, func(gtx layout.Context) layout.Dimensions {
				return drawComboRow(gtx, theme, item.label, i == this.selectedIndex)
			})
		}))
	}
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
	call := macro.Stop()
	radius := gtx.Dp(2)
	rect := image.Rectangle{Max: dims.Size}
	paint.FillShape(gtx.Ops, themes.ColorInput, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, themes.ColorBorder, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(1)),
	}.Op())
	call.Add(gtx.Ops)
	return dims
}

func drawComboRow(gtx layout.Context, theme *material.Theme, label string, selected bool) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Body1(theme, label)
		label.Color = themes.ColorText
		label.TextSize = unit.Sp(13)
		label.MaxLines = 1
		label.Truncator = "…"
		if selected {
			label.Color = themes.ColorAccent
			label.Font = font.Font{Weight: font.SemiBold}
		}
		return label.Layout(gtx)
	})
	call := macro.Stop()
	if selected {
		paint.FillShape(gtx.Ops, themes.ColorSelection,
			clip.Rect{Max: dims.Size}.Op())
	}
	call.Add(gtx.Ops)
	if dims.Size.X < gtx.Constraints.Min.X {
		dims.Size.X = gtx.Constraints.Min.X
	}
	return dims
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
