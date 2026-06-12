package components

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
)

// TODO: Move footer to preview panel

type FooterPanel struct {
	btnGenerate     widget.Clickable
	btnSaveTemplate widget.Clickable
	btnPickOutput   widget.Clickable
	btnRevealOutput widget.Clickable

	state *State
}

func NewFooterPanel(state *State) *FooterPanel {
	panel := FooterPanel{state: state}
	return &panel
}

func (this *FooterPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	return widgets.NewPanelWidget(unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(120)
						label := material.Body1(theme, "Output dir")
						label.Color = themes.ColorText
						label.TextSize = unit.Sp(13)
						return label.Layout(gtx)
					}),
					layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.state.outputPath, "Choose folder")),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(widgets.NewButtonWidget(theme, "Browse…", &this.btnPickOutput, false)),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(widgets.NewButtonWidget(theme, "Reveal", &this.btnRevealOutput, false)),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						message, isErr := this.state.GetStatus()
						if message == "" {
							message = "Ready."
						}
						col := themes.ColorText
						if isErr {
							col = themes.ColorError
						} else if this.state.GetLastTemplate() != nil {
							col = themes.ColorAccentBright
						}
						label := material.Body2(theme, message)
						label.Color = col
						label.TextSize = unit.Sp(12)
						label.MaxLines = 2
						return label.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(190)
						return widgets.NewGoldButtonWidget(theme, "⚔  Generate Template", &this.btnGenerate, false)(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(180)
						return widgets.NewGoldButtonWidget(theme, "💾  Save Template", &this.btnSaveTemplate, this.state.GetLastTemplate() == nil)(gtx)
					}),
				)
			}),
		)
	})
}

func (this *FooterPanel) HandleClicks(gtx layout.Context) {
	if this.btnGenerate.Clicked(gtx) {
		this.state.Generate()
	}
	if this.btnSaveTemplate.Clicked(gtx) {
		this.state.SaveTemplate()
	}
	if this.btnPickOutput.Clicked(gtx) {
		this.state.PickOutputDir()
	}
	if this.btnRevealOutput.Clicked(gtx) {
		_ = utils.RevealInExplorer(strings.TrimSpace(this.state.outputPath.Text()))
	}
}
