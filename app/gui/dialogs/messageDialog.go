package dialogs

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

// MessageDialog shows an informational message that closes on OK.
type MessageDialog struct {
	title string
	body  string
	okBtn widget.Clickable
}

func NewMessageDialog(title, body string) interfaces.IDialog {
	return &MessageDialog{title: title, body: body}
}

func (this *MessageDialog) Title() string { return this.title }

func (this *MessageDialog) PreferredSize() (width, height unit.Dp) { return 460, 220 }

func (this *MessageDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	closed := this.okBtn.Clicked(gtx)
	dimensions := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, widgets.NewLabelBuilder(theme).WithSizeBig().WithText(this.body).Build),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, 0)}
				}),
				layout.Rigid(widgets.NewBrightButtonWidget(theme, "OK", &this.okBtn, false)),
			)
		}),
	)
	return dimensions, closed
}
