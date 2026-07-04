package editor

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/panels"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

type Window struct {
	state *drivers.State

	tabs        []*drivers.Tab
	selectedTab int

	toolbar      *Toolbar
	previewPanel *panels.PreviewPanel
}

func NewWindow() *Window {
	window := Window{state: drivers.NewUIState()}
	window.toolbar = NewToolbar(window.state, window.load)
	window.tabs = []*drivers.Tab{
		drivers.NewTab("General", panels.NewGeneralPanel(window.state)),
		drivers.NewTab("Layout & Zones", panels.NewLayoutPanel(window.state)),
		drivers.NewTab("Bonuses & Bans", panels.NewBonusesPanel(window.state)),
	}
	window.previewPanel = panels.NewPreviewPanel(window.state)
	window.tabs[0].SetSelected(true)
	return &window
}

func (this *Window) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.save()
	if redrawAt, scheduleRedraw := this.state.AutoRegenerate(gtx.Now); scheduleRedraw {
		gtx.Execute(op.InvalidateCmd{At: redrawAt})
	}
	this.handleClicks(gtx)

	paint.FillShape(gtx.Ops, themes.ColorBackground, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(this.toolbar.GetWidget(theme)),
					layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
					layout.Rigid(this.getTabsWidget(gtx, theme)),
					layout.Flexed(1, this.getPanelsWidget(theme)))
			})
		}),
		layout.Expanded(this.state.Dialogs().GetActiveDialogWidget(theme)),
	)
}

func (this *Window) getTabsWidget(gtx layout.Context, theme *material.Theme) layout.Widget {
	this.updateTabs(gtx)

	children := make([]layout.FlexChild, 0)
	for _, tab := range this.tabs {
		children = append(children, layout.Rigid(tab.GetWidget(theme)))
	}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
	}
}

func (this *Window) getPanelsWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, widgets.NewPanelWidget(unit.Dp(0), func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(constants.DefaultPaddingLarge).
					Layout(gtx, this.getSelectedPanelWidget(theme))
			})),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(constants.DefaultPreviewWidthMinimum)
				gtx.Constraints.Max.X = gtx.Dp(constants.DefaultPreviewWidthMaximum)
				return this.previewPanel.GetPanelWidget(theme)(gtx)
			}))
	}
}

func (this *Window) getSelectedPanelWidget(theme *material.Theme) layout.Widget {
	return this.tabs[this.selectedTab].GetPanelWidget(theme)
}

func (this *Window) updateTabs(gtx layout.Context) {
	for i, tab := range this.tabs {
		if tab.IsTabClicked(gtx) && this.selectedTab != i {
			this.selectedTab = i
			this.updateTabSelection()
		}
	}
}

func (this *Window) updateTabSelection() {
	for i, tab := range this.tabs {
		tab.SetSelected(this.selectedTab == i)
	}
}

func (this *Window) handleClicks(gtx layout.Context) {
	this.toolbar.HandleClicks(gtx)
	this.previewPanel.HandleClicks(gtx)
}

func (this *Window) save() {
	for _, tab := range this.tabs {
		tab.SaveToState()
	}
}

func (this *Window) load() {
	for _, tab := range this.tabs {
		tab.LoadFromState()
	}
}
