package components

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
)

type Toolbar struct {
	buttonReset  widget.Clickable
	buttonOpen   widget.Clickable
	buttonSave   widget.Clickable
	buttonSaveAs widget.Clickable

	state *State
}

func NewToolbar(state *State) *Toolbar {
	return &Toolbar{
		state: state,
	}
}

func (this *Toolbar) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		row := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
		return row.Layout(gtx,
			layout.Rigid(widgets.NewButtonWidget(theme, "📄 New", &this.buttonReset, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "📂 Open…", &this.buttonOpen, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "💾 Save", &this.buttonSave, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "💾 Save As…", &this.buttonSaveAs, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(12)),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					path := this.state.GetCurrentPath()
					if path == "" {
						path = "(unsaved)"
					}
					if this.state.IsDirty() {
						path += " *"
					}
					label := material.Body2(theme, "File: "+path)
					label.Color = themes.ColorTextDim
					label.TextSize = unit.Sp(11)
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
		// this.seedDefaultPlayerZoneContent()
		// this.applyFromSettingsFile()
	}
	if this.buttonOpen.Clicked(gtx) {
		this.state.Load()
	}
	if this.buttonSave.Clicked(gtx) {
		this.state.Save()
	}
	if this.buttonSaveAs.Clicked(gtx) {
		this.state.SaveAs(this.state.GetSettingsFile().TemplateName)
	}
	// if this.buttonTemplates.Clicked(gtx) {
	// 	this.openTemplatesFolder()
	// }
}
