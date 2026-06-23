package editor

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

type Toolbar struct {
	buttonReset  widget.Clickable
	buttonOpen   widget.Clickable
	buttonSave   widget.Clickable
	buttonSaveAs widget.Clickable
	buttonExit   widget.Clickable

	resetCallback func()
	state         *drivers.State
}

func NewToolbar(state *drivers.State, resetCallback func()) *Toolbar {
	return &Toolbar{
		state:         state,
		resetCallback: resetCallback,
	}
}

func (this *Toolbar) GetWidget(theme *material.Theme) layout.Widget {
	// TODO: try using this, maybe it will work instead of ascii icons
	// "gioui.org/x/vector/icon"
	// "golang.org/x/exp/shiny/materialdesign/icons"
	// docIcon, _ := icon.NewIcon(icons.ActionDescription) // Standard file description vector
	// return layout.Flex{
	//     Axis:      layout.Horizontal,
	//     Alignment: layout.Middle,
	// }.Layout(gtx,
	//     // The Vector Icon
	//     layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//         return layout.Inset{Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	//             // Set explicitly matching icon color and size
	//             return docIcon.Layout(gtx, th.Palette.Fg)
	//         })
	//     }),
	//     // The Text
	//     layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	//         return material.Body2(th, "New").Layout(gtx)
	//     }),
	// )

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(widgets.NewButtonWidget(theme, "New", &this.buttonReset, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Load", &this.buttonOpen, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Save", &this.buttonSave, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Save As", &this.buttonSaveAs, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Exit", &this.buttonExit, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(12)),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					path := this.state.GetCurrentPath()
					if path == "" {
						path = "(unsaved)"
					}
					if this.state.IsUnsaved() {
						path += " *"
					}
					label := material.Caption(theme, "File: "+path)
					label.Color = themes.ColorTextDim
					label.MaxLines = 1
					label.Truncator = "…"
					label.Alignment = text.End
					return label.Layout(gtx)
				})
			}),
		)
	}
}

func (this *Toolbar) HandleClicks(gtx layout.Context) {
	if this.buttonReset.Clicked(gtx) {
		this.state.Reset()
		this.resetCallback()
	}
	if this.buttonOpen.Clicked(gtx) {
		this.state.Load()
	}
	if this.buttonSave.Clicked(gtx) {
		this.state.Save()
	}
	if this.buttonSaveAs.Clicked(gtx) {
		this.state.SaveAs(this.state.GetStateData().TemplateName)
	}
	if this.buttonExit.Clicked(gtx) {
		this.state.Exit()
	}
}
