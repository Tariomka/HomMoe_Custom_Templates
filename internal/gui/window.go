package gui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
)

type Window struct {
	state *components.State

	tabs        []*components.Tab
	selectedTab int

	toolbar      *components.Toolbar
	previewPanel *components.PreviewPanel
	footerPanel  *components.FooterPanel
}

func NewWindow() *Window {
	window := Window{state: components.NewUIState()}
	window.toolbar = components.NewToolbar(window.state, window.reset)
	window.tabs = []*components.Tab{
		components.NewTab("General", components.NewBasicSetupPanel(window.state)),
		components.NewTab("Layout", components.NewGenerationPanel(window.state)),
		components.NewTab("Game Rules (TO BE DELETED)", components.NewRulesPanel(window.state)),
		components.NewTab("Zone Content", components.NewZoneContentPanel(window.state)),
		components.NewTab("Bonuses & Bans", components.NewBonusesPanel(window.state)),
	}
	window.previewPanel = components.NewPreviewPanel(window.state)
	window.footerPanel = components.NewFooterPanel(window.state)
	window.tabs[0].SetSelected(true)
	return &window
}

func (this *Window) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.save()

	this.toolbar.HandleClicks(gtx)
	this.previewPanel.HandleClicks(gtx)
	this.footerPanel.HandleClicks(gtx)

	paint.FillShape(gtx.Ops, themes.ColorBackground, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())

	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(widgets.NewTitleBarWidget(theme, "⚔  Olden Era — Template Generator")),
			layout.Rigid(widgets.NewVerticalSpacerWidget(6)),
			layout.Rigid(this.toolbar.GetWidget(theme)),
			layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
			layout.Rigid(this.getTabsWidget(gtx, theme)),
			layout.Flexed(1, this.getPanelsWidget(theme)),
			layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
			layout.Rigid(this.footerPanel.GetPanelWidget(theme)))
	})
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
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, this.getSelectedPanelWidget(theme))
			})),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(unit.Dp(380))
				gtx.Constraints.Max.X = gtx.Dp(unit.Dp(440))
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

func (this *Window) reset() {
	this.state.Reset()
	this.load()
}
