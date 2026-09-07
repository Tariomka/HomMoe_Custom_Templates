package drivers

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
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

// DialogHost renders a stack of modal Dialogs over the main UI; only the
// top-most is shown and interactive. The stack lets a dialog open another
// dialog on top of itself (e.g. the bonus picker launching the spell picker)
// and resume once the child is dismissed. It is stored on State so any
// component can open a dialog from deep in the tree.
type DialogHost struct {
	stack    []interfaces.IDialog
	scrim    widget.Clickable
	closeBtn widget.Clickable
}

// Open pushes the given dialog onto the stack, making it the active modal.
func (this *DialogHost) Open(dialog interfaces.IDialog) { this.stack = append(this.stack, dialog) }

// Close dismisses the top-most modal, if any, resuming the one beneath it.
func (this *DialogHost) Close() {
	if len(this.stack) > 0 {
		this.stack = this.stack[:len(this.stack)-1]
	}
}

// IsOpen reports whether any modal is currently shown.
func (this *DialogHost) IsOpen() bool { return len(this.stack) > 0 }

func (this *DialogHost) GetActiveDialogWidget(theme *material.Theme) layout.Widget {
	active := this.getTopDialog()
	if active == nil {
		return func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{Size: gtx.Constraints.Min} }
	}

	return func(gtx layout.Context) layout.Dimensions {
		full := gtx.Constraints.Max

		// Scrim: a darkened, click-absorbing barrier over the entire window so the
		// underlying UI cannot be interacted with while the modal is open.
		paint.FillShape(gtx.Ops, themes.ColorsBase.Scrim, clip.Rect{Max: full}.Op())
		this.scrim.Layout(gtx, func(layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: full}
		})
		this.scrim.Clicked(gtx) // drain; clicking the scrim is intentionally a no-op.

		// Determine the panel size, clamped to the available window space.
		margin := gtx.Dp(unit.Dp(24))
		prefW, prefH := active.PreferredSize()
		panelW := gtx.Dp(unit.Dp(520))
		panelH := gtx.Dp(unit.Dp(460))
		if prefW > 0 {
			panelW = gtx.Dp(prefW)
		}
		if prefH > 0 {
			panelH = gtx.Dp(prefH)
		}
		panelW = min(panelW, full.X-margin)
		panelH = min(panelH, full.Y-margin)

		offset := image.Pt((full.X-panelW)/2, (full.Y-panelH)/2)
		defer op.Offset(offset).Push(gtx.Ops).Pop()

		gtx.Constraints = layout.Exact(image.Pt(panelW, panelH))
		return this.layoutPanel(gtx, theme)
	}
}

func (this *DialogHost) layoutPanel(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Absorb clicks that land on the panel so they do not fall through to the
	// scrim (and to keep the panel area visually solid).
	size := gtx.Constraints.Max
	radius := gtx.Dp(constants.DefaultRoundnessLarge)
	rect := image.Rectangle{Max: size}
	paint.FillShape(gtx.Ops, themes.ColorsBase.Panel, clip.UniformRRect(rect, radius).Op(gtx.Ops))
	paint.FillShape(gtx.Ops, themes.ColorsBase.Accent, clip.Stroke{
		Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
		Width: float32(gtx.Dp(unit.Dp(1))),
	}.Op())

	return layout.UniformInset(constants.DefaultPaddingLarge+4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(this.getHeaderWidget(theme)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				width := gtx.Constraints.Max.X
				height := gtx.Dp(unit.Dp(1))
				paint.FillShape(gtx.Ops, themes.ColorsBase.Border, clip.Rect{Max: image.Pt(width, height)}.Op())
				return layout.Dimensions{Size: image.Pt(width, height)}
			}),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(10)),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				active := this.getTopDialog()
				if active == nil {
					return layout.Dimensions{}
				}

				dims, done := active.Body(gtx, theme)
				if done {
					this.Close()
					gtx.Execute(op.InvalidateCmd{}) // Schedule a redraw
				}
				return dims
			}),
		)
	})
}

func (this *DialogHost) getHeaderWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if this.closeBtn.Clicked(gtx) {
			this.Close()
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				title := material.Subtitle1(theme, this.getDialogTitle())
				title.Color = themes.ColorsBase.AccentBright
				title.Font = font.Font{Weight: font.SemiBold}
				return title.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Clickable(gtx, &this.closeBtn, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(constants.DefaultPaddingSmall-2).
						Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(theme, "X")
							label.Color = themes.ColorsBase.TextDim
							if this.closeBtn.Hovered() {
								label.Color = themes.ColorsBase.AccentBright
							}
							return label.Layout(gtx)
						})
				})
			}),
		)
	}
}

// getTopDialog returns the active (top-most) dialog, or nil when the stack is empty.
func (this *DialogHost) getTopDialog() interfaces.IDialog {
	if len(this.stack) == 0 {
		return nil
	}

	return this.stack[len(this.stack)-1]
}

// getDialogTitle returns the active dialog's title, or empty when none is shown.
func (this *DialogHost) getDialogTitle() string {
	if active := this.getTopDialog(); active != nil {
		return active.Title()
	}

	return ""
}
